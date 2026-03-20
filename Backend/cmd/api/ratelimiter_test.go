package main

// func TestRatelimiter(t *testing.T) {
// 	app := newTestApplication(t)
// 	mux := app.mount()
// 	t.Run("rate limit test", func(t *testing.T) {
// 		for range 200 {
// 			req, err := http.NewRequest(http.MethodGet, "/v1/health", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			resp := executeRequest(req, mux)

// 			checkResponseCode(t, http.StatusTooManyRequests, resp.Code)
// 		}
// 	})
// }
