package server

import (
	"log"
	"net/http"
	"strconv"

	"pocket-bot/pkg/repository"

	"github.com/zhashkevych/go-pocket-sdk"
)

type AuthServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectUrl     string
}

func NewAuthServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectUrl string) *AuthServer {
	return &AuthServer{
		pocketClient:    pocketClient,
		tokenRepository: tokenRepository,
		redirectUrl:     redirectUrl,
	}
}

func (s *AuthServer) Start() error {
	s.server = &http.Server{
		Addr:    ":8080",
		Handler: s,
	}
	return s.server.ListenAndServe()
}

func (s *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatIDParam := r.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqToken, err := s.tokenRepository.Get(chatID, repository.RequestToken)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	auhResp, err := s.pocketClient.Authorize(r.Context(), reqToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.tokenRepository.Save(chatID, auhResp.AccessToken, repository.AccessToken); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("chat_id: %d\nrequest_token: %s\naccess_token: %s\n", chatID, reqToken, auhResp.AccessToken)

	w.Header().Add("Location", s.redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)

}
