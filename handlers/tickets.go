package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func TicketsHandler(c *gin.Context) {
	jiraDomain := os.Getenv("JIRA_DOMAIN")
	jiraApiToken := os.Getenv("JIRA_API_TOKEN")
	jiraEmail := os.Getenv("JIRA_EMAIL")

	ticketKey := c.Query("key")

	var jql string
	if ticketKey != "" {
		jql = fmt.Sprintf(`key = '%s'`, ticketKey)
	} else {
		jql = `project = 'Bus' AND status = 'To Do' AND assignee is EMPTY ORDER BY created DESC`
	}
	log.Println("#################################JQL Query:", jql)
	url := fmt.Sprintf("https://%s/rest/api/3/search", jiraDomain)

	payload := fmt.Sprintf(`{
        "jql": "%s", 
        "maxResults": 50, 
        "fields": ["key", "summary", "status", "assignee", "description"]
    }`, jql)
	req, err := http.NewRequest("POST", url, io.NopCloser(strings.NewReader(payload)))
	if err != nil {
		log.Println("Erro ao criar request:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(jiraEmail, jiraApiToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Erro ao enviar request:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar request"})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		log.Printf("Erro da API Jira (%d): %s", res.StatusCode, body)
		c.JSON(res.StatusCode, gin.H{"error": fmt.Sprintf("Erro na resposta da API Jira: %s", string(body))})
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Printf("Erro ao parsear JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados"})
		return
	}

	var issues []map[string]interface{}
	if issuesData, exists := result["issues"].([]interface{}); exists {
		for _, issueData := range issuesData {
			issue := issueData.(map[string]interface{})
			fields := issue["fields"].(map[string]interface{})
			description := fields["description"]
			plainText := extractTextFromADF(description)

			issues = append(issues, map[string]interface{}{
				"key":         issue["key"],
				"summary":     fields["summary"],
				"status":      fields["status"],
				"description": plainText,
			})
		}
	}

	if ticketKey != "" {
		if len(issues) > 0 {
			c.JSON(http.StatusOK, issues[0])
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"issues": issues})
}

func extractTextFromADF(adf interface{}) string {
	var text []string

	switch v := adf.(type) {
	case map[string]interface{}:

		if content, exists := v["content"]; exists {
			for _, item := range content.([]interface{}) {
				text = append(text, extractTextFromADF(item))
			}
		}

		if t, exists := v["text"]; exists {
			text = append(text, t.(string))
		}
	}

	return strings.Join(text, " ")
}
