# Solution to Exercise 5.2

This is my workflow in creating the solution:

1. Copy all the code from (`chap5/streaming-decode`)
2. Update the go module name to: `github.com/practicalgo/book-exercise-solutions/chap5/exercise2`
3. We will create a new test function to verify the behavior we will implement in this exericse:

```
func Test_DecodeUnknownFieldError(t *testing.T) {
	const jsonStream = `
	{"user_ip": "172.121.19.21", "event": "click_on_add_cart", "user_data":"some_data"}{"user_ip": "172.121.19.21", "event": "click_on_checkout"}
`
	body := strings.NewReader(jsonStream)

	r := httptest.NewRequest("POST", "http://example.com/decode", body)
	w := httptest.NewRecorder()

	decodeHandler(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected Response Status: %v, Got: %v", http.StatusBadRequest, w.Result().StatusCode)
	}
}
```

Run the test, and see it fail.

4. Update `decodeHandler()` function to call the method `DisallowUnknownFields()` after
creating the decoder:

```
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields()
```

The test above should now pass.