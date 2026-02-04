package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/resume-builder/backend/internal/config"
	"github.com/resume-builder/backend/internal/models"
	"github.com/rs/zerolog/log"
)

type AIService struct {
	client *anthropic.Client
	model  string
}

func NewAIService(cfg *config.Config) *AIService {
	client := anthropic.NewClient(
		option.WithAPIKey(cfg.AnthropicAPIKey),
	)

	return &AIService{
		client: client,
		model:  cfg.AnthropicModel,
	}
}

// ExtractionResult contains extracted resume data with confidence scores
type ExtractionResult struct {
	Data       *models.ResumeData    `json:"data"`
	Confidence map[string]float64    `json:"confidence"`
	Warnings   []string              `json:"warnings,omitempty"`
}

// ExtractResumeData extracts structured resume data from raw text
func (s *AIService) ExtractResumeData(ctx context.Context, rawText string) (*ExtractionResult, error) {
	prompt := buildExtractionPrompt(rawText)

	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(s.model),
		MaxTokens: anthropic.F(int64(4000)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(extractionSystemPrompt),
		}),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to call Anthropic API for extraction")
		return nil, fmt.Errorf("AI extraction failed: %w", err)
	}

	// Parse the response
	var responseText string
	for _, block := range message.Content {
		if block.Type == anthropic.ContentBlockTypeText {
			responseText = block.Text
			break
		}
	}

	// Parse JSON from response
	result, err := parseExtractionResponse(responseText)
	if err != nil {
		log.Error().Err(err).Str("response", responseText).Msg("Failed to parse extraction response")
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return result, nil
}

// ImprovementVariant represents one variant of improved text
type ImprovementVariant struct {
	Type        string `json:"type"` // concise, detailed, ats_optimized
	Text        string `json:"text"`
	Explanation string `json:"explanation"`
}

// ImprovementResult contains multiple variants of improved text
type ImprovementResult struct {
	Original string               `json:"original"`
	Variants []ImprovementVariant `json:"variants"`
}

// ImproveText provides AI-powered text improvement suggestions
func (s *AIService) ImproveText(ctx context.Context, text string, context string) (*ImprovementResult, error) {
	prompt := buildImprovementPrompt(text, context)

	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(s.model),
		MaxTokens: anthropic.F(int64(1500)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(improvementSystemPrompt),
		}),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to call Anthropic API for improvement")
		return nil, fmt.Errorf("AI improvement failed: %w", err)
	}

	var responseText string
	for _, block := range message.Content {
		if block.Type == anthropic.ContentBlockTypeText {
			responseText = block.Text
			break
		}
	}

	result, err := parseImprovementResponse(responseText, text)
	if err != nil {
		log.Error().Err(err).Str("response", responseText).Msg("Failed to parse improvement response")
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return result, nil
}

// GenerateSummary creates a professional summary based on resume data
func (s *AIService) GenerateSummary(ctx context.Context, data *models.ResumeData) (string, error) {
	prompt := buildSummaryPrompt(data)

	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(s.model),
		MaxTokens: anthropic.F(int64(300)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(summarySystemPrompt),
		}),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to call Anthropic API for summary generation")
		return "", fmt.Errorf("AI summary generation failed: %w", err)
	}

	var responseText string
	for _, block := range message.Content {
		if block.Type == anthropic.ContentBlockTypeText {
			responseText = block.Text
			break
		}
	}

	return responseText, nil
}

// JobAnalysisResult contains the analysis of a resume against a job description
type JobAnalysisResult struct {
	MatchScore         int      `json:"match_score"`
	MatchingSkills     []string `json:"matching_skills"`
	MissingSkills      []string `json:"missing_skills"`
	Recommendations    []string `json:"recommendations"`
	KeywordsToAdd      []string `json:"keywords_to_add"`
	SectionsToEmphasize []string `json:"sections_to_emphasize"`
}

// AnalyzeJobMatch analyzes how well a resume matches a job description
func (s *AIService) AnalyzeJobMatch(ctx context.Context, resumeData *models.ResumeData, jobDescription string) (*JobAnalysisResult, error) {
	prompt := buildJobAnalysisPrompt(resumeData, jobDescription)

	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(s.model),
		MaxTokens: anthropic.F(int64(2000)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(jobAnalysisSystemPrompt),
		}),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to call Anthropic API for job analysis")
		return nil, fmt.Errorf("AI job analysis failed: %w", err)
	}

	var responseText string
	for _, block := range message.Content {
		if block.Type == anthropic.ContentBlockTypeText {
			responseText = block.Text
			break
		}
	}

	result, err := parseJobAnalysisResponse(responseText)
	if err != nil {
		log.Error().Err(err).Str("response", responseText).Msg("Failed to parse job analysis response")
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return result, nil
}

// System prompts
const extractionSystemPrompt = `You are an expert resume parser. Your task is to extract structured information from resume text.

Always respond with valid JSON in the following format:
{
  "data": {
    "personal_info": {
      "first_name": "",
      "last_name": "",
      "title": "",
      "email": "",
      "phone": "",
      "location": "",
      "website": "",
      "linkedin": "",
      "github": ""
    },
    "summary": "",
    "experience": [
      {
        "id": "unique-id",
        "company": "",
        "position": "",
        "location": "",
        "start_date": "YYYY-MM",
        "end_date": "YYYY-MM or empty if current",
        "is_current": false,
        "description": "",
        "achievements": [],
        "technologies": [],
        "order": 0
      }
    ],
    "education": [
      {
        "id": "unique-id",
        "institution": "",
        "degree": "",
        "field_of_study": "",
        "location": "",
        "start_date": "YYYY-MM",
        "end_date": "YYYY-MM",
        "gpa": "",
        "honors": "",
        "order": 0
      }
    ],
    "skills": {
      "technical": [{"name": "", "level": ""}],
      "soft": [],
      "languages": [{"name": "", "proficiency": ""}]
    },
    "certifications": [],
    "projects": []
  },
  "confidence": {
    "personal_info": 0.95,
    "experience": 0.90,
    "education": 0.85,
    "skills": 0.80
  },
  "warnings": []
}

Rules:
- Extract all available information
- Use null for missing fields
- Provide confidence scores (0.0-1.0) for each section
- Add warnings for ambiguous or potentially incorrect data
- Format dates as YYYY-MM when possible
- Generate unique IDs for each item`

const improvementSystemPrompt = `You are an expert resume writer. Your task is to improve resume content to be more impactful and professional.

Always respond with valid JSON in the following format:
{
  "variants": [
    {
      "type": "concise",
      "text": "Improved text that is shorter and more direct",
      "explanation": "Why this version is effective"
    },
    {
      "type": "detailed",
      "text": "Improved text with more specific details and metrics",
      "explanation": "Why this version is effective"
    },
    {
      "type": "ats_optimized",
      "text": "Version optimized for ATS systems with relevant keywords",
      "explanation": "Why this version is effective"
    }
  ]
}

Guidelines:
- Use strong action verbs (Led, Developed, Implemented, Achieved)
- Include quantifiable metrics when possible
- Avoid clichés and buzzwords
- Be specific about accomplishments
- Use industry-appropriate terminology
- Keep ATS version simple and keyword-rich`

const summarySystemPrompt = `You are an expert resume writer. Create a compelling professional summary based on the provided career information.

Guidelines:
- Keep it to 2-3 sentences
- Highlight key strengths and experience
- Mention years of experience if significant
- Include relevant skills and achievements
- Use confident, professional tone
- Avoid generic statements

Respond with only the summary text, no JSON formatting.`

const jobAnalysisSystemPrompt = `You are an expert career advisor. Analyze how well a resume matches a job description.

Always respond with valid JSON:
{
  "match_score": 75,
  "matching_skills": ["skill1", "skill2"],
  "missing_skills": ["skill3", "skill4"],
  "recommendations": ["suggestion1", "suggestion2"],
  "keywords_to_add": ["keyword1", "keyword2"],
  "sections_to_emphasize": ["Experience", "Skills"]
}

Analysis criteria:
- Required skills match
- Preferred qualifications match
- Experience level alignment
- Industry terminology usage
- Keyword presence`

func buildExtractionPrompt(rawText string) string {
	return fmt.Sprintf(`Extract structured resume data from the following text:

---
%s
---

Analyze the text carefully and extract all relevant information into the structured JSON format.`, rawText)
}

func buildImprovementPrompt(text, context string) string {
	contextInfo := ""
	if context != "" {
		contextInfo = fmt.Sprintf("\nContext: %s", context)
	}
	return fmt.Sprintf(`Improve the following resume text:%s

Original text:
"%s"

Provide three improved variants (concise, detailed, and ATS-optimized).`, contextInfo, text)
}

func buildSummaryPrompt(data *models.ResumeData) string {
	// Build context from resume data
	experienceYears := len(data.Experience)
	var skills []string
	for _, s := range data.Skills.Technical {
		skills = append(skills, s.Name)
	}
	
	return fmt.Sprintf(`Create a professional summary for someone with the following background:

Title: %s
Experience: %d positions
Key Skills: %v

Most recent role: %s at %s

Generate a compelling 2-3 sentence professional summary.`,
		data.PersonalInfo.Title,
		experienceYears,
		skills,
		getLatestPosition(data),
		getLatestCompany(data))
}

func buildJobAnalysisPrompt(resumeData *models.ResumeData, jobDescription string) string {
	resumeJSON, _ := json.Marshal(resumeData)
	return fmt.Sprintf(`Analyze how well this resume matches the job description.

Resume Data:
%s

Job Description:
%s

Provide a detailed analysis with match score, matching/missing skills, and recommendations.`, string(resumeJSON), jobDescription)
}

func getLatestPosition(data *models.ResumeData) string {
	if len(data.Experience) > 0 {
		return data.Experience[0].Position
	}
	return "N/A"
}

func getLatestCompany(data *models.ResumeData) string {
	if len(data.Experience) > 0 {
		return data.Experience[0].Company
	}
	return "N/A"
}

func parseExtractionResponse(response string) (*ExtractionResult, error) {
	// Find JSON in response (it might be wrapped in markdown code blocks)
	jsonStr := extractJSON(response)
	
	var result ExtractionResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func parseImprovementResponse(response string, original string) (*ImprovementResult, error) {
	jsonStr := extractJSON(response)
	
	var parsed struct {
		Variants []ImprovementVariant `json:"variants"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, err
	}
	
	return &ImprovementResult{
		Original: original,
		Variants: parsed.Variants,
	}, nil
}

func parseJobAnalysisResponse(response string) (*JobAnalysisResult, error) {
	jsonStr := extractJSON(response)
	
	var result JobAnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func extractJSON(text string) string {
	// Simple extraction - find content between first { and last }
	start := -1
	end := -1
	depth := 0
	
	for i, c := range text {
		if c == '{' {
			if start == -1 {
				start = i
			}
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 {
				end = i + 1
				break
			}
		}
	}
	
	if start != -1 && end != -1 {
		return text[start:end]
	}
	return text
}
