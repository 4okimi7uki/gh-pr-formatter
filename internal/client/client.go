package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/4okimi7uki/gh-pr-formatter/internal/models"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
	token      string
}

type pullRequestsResponse struct {
	Data struct {
		Repository struct {
			Main struct {
				Nodes []struct {
					Number      int       `json:"number"`
					HeadRefName string    `json:"headRefName"`
					MergedAt    time.Time `json:"mergedAt"`
				} `json:"nodes"`
			} `json:"main"`
			Develop struct {
				Nodes []struct {
					Number   int       `json:"number"`
					Title    string    `json:"title"`
					MergedAt time.Time `json:"mergedAt"`
					Author   struct {
						Login string `json:"login"`
					} `json:"author"`
				} `json:"nodes"`
			} `json:"develop"`
		} `json:"repository"`
	} `json:"data"`
}

const githubGraphQLEndpoint = "https://api.github.com/graphql"

func NewClient(token string) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		endpoint:   githubGraphQLEndpoint,
		token:      token,
	}
}

func (c *Client) Do(query string, vars map[string]any, v any) error {
	body := map[string]any{
		"query":     query,
		"variables": vars,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub GraphQL error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func filterPullRequestsSince(from time.Time, prs []models.MergedPrList) []models.MergedPrList {
	var filteredPrs []models.MergedPrList
	for _, pr := range prs {
		if pr.MergedAt.After(from) {
			filteredPrs = append(filteredPrs, pr)
		}
	}

	return filteredPrs
}

func (c *Client) ListMergedPullRequests(
	owner,
	name string,
	developLimit int,
) ([]models.MergedPrList, error) {
	const mainLimit = 20

	const query = `
	query($owner: String!, $name: String!, $lastMain: Int!, $lastDevelop: Int!) {
  repository(owner: $owner, name: $name) {
    main: pullRequests(
      states: MERGED
      baseRefName: "main"
      last: $lastMain
      orderBy: {field: UPDATED_AT, direction: ASC}
    ) {
      nodes {
        number
        mergedAt
        headRefName
      }
    }

    develop: pullRequests(
      states: MERGED
      baseRefName: "develop"
      last: $lastDevelop
      orderBy: {field: UPDATED_AT, direction: ASC}
    ) {
      nodes {
        number
        mergedAt
        title
        author { login }
      }
    }
  }
}`
	vars := map[string]any{
		"owner":       owner,
		"name":        name,
		"lastMain":    mainLimit,
		"lastDevelop": developLimit,
	}

	var result pullRequestsResponse

	if err := c.Do(query, vars, &result); err != nil {
		return nil, err
	}

	var mainPrs []models.MergedPrList
	mainNodes := result.Data.Repository.Main.Nodes
	for _, node := range mainNodes {
		if strings.HasPrefix(node.HeadRefName, "hotfix/") {
			continue
		}
		mainPrs = append(mainPrs, models.MergedPrList{
			PrList: models.PrList{
				Number:   node.Number,
				MergedAt: node.MergedAt,
			},
		})
	}

	var devPrs []models.MergedPrList
	developNodes := result.Data.Repository.Develop.Nodes
	for _, node := range developNodes {
		devPrs = append(devPrs, models.MergedPrList{
			PrList: models.PrList{
				Number:   node.Number,
				MergedAt: node.MergedAt,
			},
			Title: node.Title,
			Author: models.Author{
				Login: node.Author.Login,
			},
		})
	}

	if len(mainPrs) == 0 {
		return nil, fmt.Errorf("no non-hotfix merged pull requests found on main")
	}

	from := mainPrs[len(mainPrs)-1].MergedAt

	filtered := filterPullRequestsSince(from, devPrs)

	return filtered, nil
}
