package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"bailanysta/api/internal/pkg/logger"
	"bailanysta/api/internal/services"
)

type AIHandler struct {
	aiService *services.AIService
	logger    *logger.Logger
	validator *validator.Validate
}

func NewAIHandler(aiService *services.AIService, logger *logger.Logger) *AIHandler {
	return &AIHandler{
		aiService: aiService,
		logger:    logger,
		validator: validator.New(),
	}
}

func (h *AIHandler) GenerateText(w http.ResponseWriter, r *http.Request) {
	var req services.GenerateTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode generate text request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Generate text validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.aiService.GenerateText(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to generate text", map[string]interface{}{
			"error":  err.Error(),
			"prompt": req.Prompt,
		})
		h.respondWithError(w, "Failed to generate text: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Text generated successfully", map[string]interface{}{
		"prompt_length":   len(req.Prompt),
		"response_length": len(response.Text),
	})

	h.respondWithJSON(w, response, http.StatusOK)
}

func (h *AIHandler) GeneratePost(w http.ResponseWriter, r *http.Request) {
	var req services.GeneratePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode generate post request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Generate post validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.aiService.GeneratePost(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to generate post", map[string]interface{}{
			"error": err.Error(),
			"topic": req.Topic,
		})
		h.respondWithError(w, "Failed to generate post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Post generated successfully", map[string]interface{}{
		"topic":           req.Topic,
		"course":          req.Course,
		"response_length": len(response.Text),
	})

	h.respondWithJSON(w, response, http.StatusOK)
}

func (h *AIHandler) GenerateComment(w http.ResponseWriter, r *http.Request) {
	var req services.GenerateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode generate comment request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Generate comment validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.aiService.GenerateComment(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to generate comment", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Failed to generate comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Comment generated successfully", map[string]interface{}{
		"response_length": len(response.Text),
	})

	h.respondWithJSON(w, response, http.StatusOK)
}

func (h *AIHandler) GenerateStudyNotes(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Topic  string `json:"topic" validate:"required,min=3,max=200"`
		Course string `json:"course,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode generate study notes request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Generate study notes validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.aiService.GenerateStudyNotes(r.Context(), req.Topic, req.Course)
	if err != nil {
		h.logger.Error("Failed to generate study notes", map[string]interface{}{
			"error": err.Error(),
			"topic": req.Topic,
		})
		h.respondWithError(w, "Failed to generate study notes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Study notes generated successfully", map[string]interface{}{
		"topic":  req.Topic,
		"course": req.Course,
	})

	h.respondWithJSON(w, response, http.StatusOK)
}

func (h *AIHandler) GenerateQuiz(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Topic  string `json:"topic" validate:"required,min=3,max=200"`
		Course string `json:"course,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode generate quiz request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Generate quiz validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.aiService.GenerateQuiz(r.Context(), req.Topic, req.Course)
	if err != nil {
		h.logger.Error("Failed to generate quiz", map[string]interface{}{
			"error": err.Error(),
			"topic": req.Topic,
		})
		h.respondWithError(w, "Failed to generate quiz: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Quiz generated successfully", map[string]interface{}{
		"topic":  req.Topic,
		"course": req.Course,
	})

	h.respondWithJSON(w, response, http.StatusOK)
}

func (h *AIHandler) ExplainConcept(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Concept string `json:"concept" validate:"required,min=3,max=200"`
		Context string `json:"context,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode explain concept request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Explain concept validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondWithError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.aiService.ExplainConcept(r.Context(), req.Concept, req.Context)
	if err != nil {
		h.logger.Error("Failed to explain concept", map[string]interface{}{
			"error":   err.Error(),
			"concept": req.Concept,
		})
		h.respondWithError(w, "Failed to explain concept: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Concept explained successfully", map[string]interface{}{
		"concept": req.Concept,
		"context": req.Context,
	})

	h.respondWithJSON(w, response, http.StatusOK)
}

func (h *AIHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *AIHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	h.respondWithJSON(w, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    getErrorCode(statusCode),
			"message": message,
		},
	}, statusCode)
}
