package interactionrouter_test

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/slack-go/slack"

	routererrors "github.com/genkami/go-slack-event-router/errors"
	ir "github.com/genkami/go-slack-event-router/interactionrouter"
	"github.com/genkami/go-slack-event-router/signature"
)

var _ = Describe("InteractionRouter", func() {
	Describe("Type", func() {
		var (
			numHandlerCalled int
			innerHandler     = ir.HandlerFunc(func(_ *slack.InteractionCallback) error {
				numHandlerCalled++
				return nil
			})
		)
		BeforeEach(func() {
			numHandlerCalled = 0
		})

		Context("when the type of the interaction callback matches to the predicate's", func() {
			It("calls the inner handler", func() {
				h := ir.Type(slack.InteractionTypeBlockActions).Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type: slack.InteractionTypeBlockActions,
				}
				err := h.HandleInteraction(callback)
				Expect(err).NotTo(HaveOccurred())
				Expect(numHandlerCalled).To(Equal(1))
			})
		})

		Context("when the type of the interaction callback differs from the predicate's", func() {
			It("calls the inner handler", func() {
				h := ir.Type(slack.InteractionTypeBlockActions).Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type: slack.InteractionTypeViewSubmission,
				}
				err := h.HandleInteraction(callback)
				Expect(err).To(Equal(routererrors.NotInterested))
				Expect(numHandlerCalled).To(Equal(0))
			})
		})
	})

	Describe("BlockAction", func() {
		var (
			numHandlerCalled int
			innerHandler     = ir.HandlerFunc(func(_ *slack.InteractionCallback) error {
				numHandlerCalled++
				return nil
			})
		)
		BeforeEach(func() {
			numHandlerCalled = 0
		})

		Context("when the interaction callback has the block_action specified by the predicate", func() {
			It("calls the inner handler", func() {
				h := ir.BlockAction("BLOCK_ID", "ACTION_ID").Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type: slack.InteractionTypeBlockActions,
					ActionCallback: slack.ActionCallbacks{
						BlockActions: []*slack.BlockAction{
							{BlockID: "BLOCK_ID", ActionID: "ACTION_ID"},
						},
					},
				}
				err := h.HandleInteraction(callback)
				Expect(err).NotTo(HaveOccurred())
				Expect(numHandlerCalled).To(Equal(1))
			})
		})

		Context("when one of the block_acsions that the interaction callback has is the one specified by the predicate", func() {
			It("calls the inner handler", func() {
				h := ir.BlockAction("BLOCK_ID", "ACTION_ID").Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type: slack.InteractionTypeBlockActions,
					ActionCallback: slack.ActionCallbacks{
						BlockActions: []*slack.BlockAction{
							{BlockID: "ANOTHER_BLOCK_ID", ActionID: "ANOTHER_ACTION_ID"},
							{BlockID: "BLOCK_ID", ActionID: "ACTION_ID"},
						},
					},
				}
				err := h.HandleInteraction(callback)
				Expect(err).NotTo(HaveOccurred())
				Expect(numHandlerCalled).To(Equal(1))
			})
		})

		Context("when the interaction callback does not have any block_action", func() {
			It("does not call the inner handler", func() {
				h := ir.BlockAction("BLOCK_ID", "ACTION_ID").Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type: slack.InteractionTypeBlockActions,
					ActionCallback: slack.ActionCallbacks{
						BlockActions: []*slack.BlockAction{},
					},
				}
				err := h.HandleInteraction(callback)
				Expect(err).To(Equal(routererrors.NotInterested))
				Expect(numHandlerCalled).To(Equal(0))
			})
		})

		Context("when the block_action in the interaction callback is not what the predicate expects", func() {
			It("does not call the inner handler", func() {
				h := ir.BlockAction("BLOCK_ID", "ACTION_ID").Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type: slack.InteractionTypeBlockActions,
					ActionCallback: slack.ActionCallbacks{
						BlockActions: []*slack.BlockAction{
							{BlockID: "ANOTHER_BLOCK_ID", ActionID: "ANOTHER_ACTION_ID"},
						},
					},
				}
				err := h.HandleInteraction(callback)
				Expect(err).To(Equal(routererrors.NotInterested))
				Expect(numHandlerCalled).To(Equal(0))
			})
		})

		Context("when the block_id in the block_action is the same as the predicate expected but the action_id isn't", func() {
			It("does not call the inner handler", func() {
				h := ir.BlockAction("BLOCK_ID", "ACTION_ID").Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type: slack.InteractionTypeBlockActions,
					ActionCallback: slack.ActionCallbacks{
						BlockActions: []*slack.BlockAction{
							{BlockID: "BLOCK_ID", ActionID: "ANOTHER_ACTION_ID"},
						},
					},
				}
				err := h.HandleInteraction(callback)
				Expect(err).To(Equal(routererrors.NotInterested))
				Expect(numHandlerCalled).To(Equal(0))
			})
		})

		Context("when the action_id in the block_action is the same as the predicate expected but the block_id isn't", func() {
			It("does not call the inner handler", func() {
				h := ir.BlockAction("BLOCK_ID", "ACTION_ID").Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type: slack.InteractionTypeBlockActions,
					ActionCallback: slack.ActionCallbacks{
						BlockActions: []*slack.BlockAction{
							{BlockID: "ANOTHER_BLOCK_ID", ActionID: "ACTION_ID"},
						},
					},
				}
				err := h.HandleInteraction(callback)
				Expect(err).To(Equal(routererrors.NotInterested))
				Expect(numHandlerCalled).To(Equal(0))
			})
		})
	})

	Describe("CallbackID", func() {
		var (
			numHandlerCalled int
			innerHandler     = ir.HandlerFunc(func(_ *slack.InteractionCallback) error {
				numHandlerCalled++
				return nil
			})
		)
		BeforeEach(func() {
			numHandlerCalled = 0
		})

		Context("when the callback_id in the interaction callback matches to the predicate's", func() {
			It("calls the inner handler", func() {
				h := ir.CallbackID("CALLBACK_ID").Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type:       slack.InteractionTypeBlockActions,
					CallbackID: "CALLBACK_ID",
				}
				err := h.HandleInteraction(callback)
				Expect(err).NotTo(HaveOccurred())
				Expect(numHandlerCalled).To(Equal(1))
			})
		})

		Context("when the callback_id in the interaction callback differs from the predicate's", func() {
			It("does not call the inner handler", func() {
				h := ir.CallbackID("CALLBACK_ID").Wrap(innerHandler)
				callback := &slack.InteractionCallback{
					Type:       slack.InteractionTypeBlockActions,
					CallbackID: "ANOTHER_CALLBACK_ID",
				}
				err := h.HandleInteraction(callback)
				Expect(err).To(Equal(routererrors.NotInterested))
				Expect(numHandlerCalled).To(Equal(0))
			})
		})
	})

	Describe("New", func() {
		Context("when neither WithSigningToken nor InsecureSkipVerification is given", func() {
			It("returns an error", func() {
				_, err := ir.New()
				Expect(err).To(MatchError(MatchRegexp("WithSigningToken")))
			})
		})

		Context("when InsecureSkipVerification is given", func() {
			It("returns a new Router", func() {
				r, err := ir.New(ir.InsecureSkipVerification())
				Expect(err).NotTo(HaveOccurred())
				Expect(r).NotTo(BeNil())
			})
		})

		Context("when WithSigningToken is given", func() {
			It("returns a new Router", func() {
				r, err := ir.New(ir.WithSigningToken("THE_TOKEN"))
				Expect(err).NotTo(HaveOccurred())
				Expect(r).NotTo(BeNil())
			})
		})

		Context("when both WithSigningToken and InsecureSkipVerification are given", func() {
			It("returns an error", func() {
				_, err := ir.New(ir.InsecureSkipVerification(), ir.WithSigningToken("THE_TOKEN"))
				Expect(err).To(MatchError(MatchRegexp("WithSigningToken")))
			})
		})
	})

	Describe("WithSigningSecret", func() {
		var (
			r       *ir.Router
			token   = "THE_TOKEN"
			content = `
			{
				"type": "shortcut",
				"token": "XXXXXXXXXXXXX",
				"action_ts": "1581106241.371594",
				"team": {
				  "id": "TXXXXXXXX",
				  "domain": "shortcuts-test"
				},
				"user": {
				  "id": "UXXXXXXXXX",
				  "username": "aman",
				  "team_id": "TXXXXXXXX"
				},
				"callback_id": "shortcut_create_task",
				"trigger_id": "944799105734.773906753841.38b5894552bdd4a780554ee59d1f3638"
			}`
		)
		BeforeEach(func() {
			var err error
			r, err = ir.New(ir.WithSigningToken(token), ir.VerboseResponse())
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the signature is valid", func() {
			It("responds with 200", func() {
				req, err := NewSignedRequest(token, content, nil)
				Expect(err).NotTo(HaveOccurred())
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("when the signature is invalid", func() {
			It("responds with Unauthorized", func() {
				req, err := NewSignedRequest(token, content, nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set(signature.HeaderSignature, "v0="+hex.EncodeToString([]byte("INVALID_SIGNATURE")))
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("when the timestamp is too old", func() {
			It("responds with Unauthorized", func() {
				ts := time.Now().Add(-1 * time.Hour)
				req, err := NewSignedRequest(token, content, &ts)
				Expect(err).NotTo(HaveOccurred())
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("InsecureSkipVerification", func() {
		var (
			r       *ir.Router
			token   = "THE_TOKEN"
			content = `
			{
				"type": "shortcut",
				"token": "XXXXXXXXXXXXX",
				"action_ts": "1581106241.371594",
				"team": {
				  "id": "TXXXXXXXX",
				  "domain": "shortcuts-test"
				},
				"user": {
				  "id": "UXXXXXXXXX",
				  "username": "aman",
				  "team_id": "TXXXXXXXX"
				},
				"callback_id": "shortcut_create_task",
				"trigger_id": "944799105734.773906753841.38b5894552bdd4a780554ee59d1f3638"
			}`
		)
		BeforeEach(func() {
			var err error
			r, err = ir.New(ir.InsecureSkipVerification(), ir.VerboseResponse())
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the signature is valid", func() {
			It("responds with 200", func() {
				req, err := NewSignedRequest(token, content, nil)
				Expect(err).NotTo(HaveOccurred())
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("when the signature is invalid", func() {
			It("responds with 200", func() {
				req, err := NewSignedRequest(token, content, nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set(signature.HeaderSignature, "v0="+hex.EncodeToString([]byte("INVALID_SIGNATURE")))
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("when the timestamp is too old", func() {
			It("responds with 200", func() {
				ts := time.Now().Add(-1 * time.Hour)
				req, err := NewSignedRequest(token, content, &ts)
				Expect(err).NotTo(HaveOccurred())
				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})
})

func NewSignedRequest(signingSecret string, payload string, ts *time.Time) (*http.Request, error) {
	var now time.Time
	if ts == nil {
		now = time.Now()
	} else {
		now = *ts
	}
	form := url.Values{}
	form.Set("payload", payload)
	body := form.Encode()
	req, err := http.NewRequest(http.MethodPost, "http://example.com/path/to/callback", bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := signature.AddSignature(req.Header, []byte(signingSecret), []byte(body), now); err != nil {
		return nil, err
	}
	return req, nil
}
