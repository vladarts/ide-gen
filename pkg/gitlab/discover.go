package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"github.com/xxxbobrxxx/ide-gen/pkg/repository"
	"os"
	"regexp"
)

type DiscoveryConfig struct {
	Url         string `json:"url,omitempty"`
	Token       string `json:"token,omitempty" jsonschema:"oneof_required=token"`
	TokenEnvVar string `json:"tokenEnvVar,omitempty" jsonschema:"oneof_required=tokenEnvVar"`
	TokenType   string `json:"tokenType,omitempty" jsonschema:"enum=private,enum=job,enum=oauth"`

	IncludeArchived bool `json:"includeArchived,omitempty"`
	HttpsUrl        bool `json:"httpsUrl,omitempty"`
	FastForward     bool `json:"fastForward,omitempty"`

	Include *[]string `json:"include,omitempty"`
	Exclude *[]string `json:"exclude,omitempty"`

	client *gitlab.Client
}

func (d *DiscoveryConfig) Discover() ([]*repository.GitSourcesRootConfig, error) {
	allProjects, err := d.ListProjects()
	if err != nil {
		return nil, err
	}

	var result []*repository.GitSourcesRootConfig

	for _, p := range allProjects {
		if p.Archived && !d.IncludeArchived {
			continue
		}

		if d.Exclude != nil {
			excluded := false
			for _, pattern := range *d.Exclude {
				matched, err := regexp.MatchString(pattern, p.PathWithNamespace)
				if err != nil {
					return nil, err
				}

				if matched {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}
		}

		if d.Include != nil {
			included := false
			for _, pattern := range *d.Include {
				matched, err := regexp.MatchString(pattern, p.PathWithNamespace)
				if err != nil {
					return nil, err
				}

				if matched {
					included = true
					break
				}
			}
			if !included {
				continue
			}
		}

		var url string
		if d.HttpsUrl {
			url = p.HTTPURLToRepo
		} else {
			url = p.SSHURLToRepo
		}

		result = append(result, &repository.GitSourcesRootConfig{
			Url:         url,
			FastForward: &d.FastForward,
		})
	}

	return result, nil
}

func (d DiscoveryConfig) ListProjects() ([]gitlab.Project, error) {
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}

	var projects []gitlab.Project

	for {
		// Get the first page with projects.
		ps, resp, err := d.client.Projects.ListProjects(opt)
		if err != nil {
			return nil, err
		}

		// List all the projects we've found so far.
		for _, p := range ps {
			projects = append(projects, *p)
		}

		// Exit the loop when we've seen all pages.
		if resp.NextPage == 0 {
			break
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	return projects, nil
}

func (d *DiscoveryConfig) Init() (err error) {
	var opts []gitlab.ClientOptionFunc
	if d.Url != "" {
		opts = append(opts, gitlab.WithBaseURL(d.Url))
	}

	var token string
	if d.TokenEnvVar != "" {
		token = os.Getenv(d.TokenEnvVar)
	} else {
		token = d.Token
	}

	var client *gitlab.Client
	switch d.TokenType {
	case "", "private":
		client, err = gitlab.NewClient(token, opts...)
	case "job":
		client, err = gitlab.NewJobClient(token, opts...)
	case "oauth":
		client, err = gitlab.NewOAuthClient(token, opts...)
	default:
		err = fmt.Errorf("incorrect token type: %v", d.TokenType)
	}
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	d.client = client
	return nil
}
