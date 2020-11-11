package fcm

import (
	"fmt"

	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Context("handleResponse func", func() {
		var (
			fuzzer = fuzz.New().NilChance(0)

			statusCode  int
			resp        sendResponse
			respBody    []byte
			expectedErr error
		)

		When("status code is success - 2XX", func() {
			It("should succeed", func() {
				statusCode = 200
				err := handleResponse(statusCode, respBody)
				Ω(err).Should(Succeed())
			})
		})

		When("status code not success != 2XX", func() {
			BeforeEach(func() {
				statusCode = 403
			})

			When("error field is missing in the sendResponse", func() {
				It("should succeed", func() {
					respBody = []byte(`{"field":"value"}`)
					expectedErr = fmt.Errorf(`empty error in sendResponse for status code 403: %s`, string(respBody))

					err := handleResponse(statusCode, respBody)
					Ω(err).Should(Equal(expectedErr))
				})
			})

			When("sendResponse body is like expected format", func() {
				BeforeEach(func() {
					fuzzer.Fuzz(&resp)
				})

				JustBeforeEach(func() {
					data, err := resp.MarshalJSON()
					Ω(err).ShouldNot(HaveOccurred())

					respBody = data
				})

				When("errorDetails are missing in sendResponse", func() {
					BeforeEach(func() {
						resp.Error.Details = nil
					})

					It("should return general error", func() {
						expectedErr = fmt.Errorf("unsuccessful sendResponse with status code: %d: %v",
							statusCode, string(respBody))

						err := handleResponse(statusCode, respBody)
						Ω(err).Should(Equal(expectedErr))
					})
				})

				When("errorDetails contain unregistered error code", func() {
					BeforeEach(func() {
						resp.Error.Details = []errorDetails{
							{
								ErrorCode: errorCodeUnregistered,
							},
						}
					})

					It("should return unregistered error", func() {
						err := handleResponse(statusCode, respBody)
						Ω(err).Should(Equal(ErrUnregistered))
					})
				})

				When("errorDetails contain unspecified error code", func() {
					BeforeEach(func() {
						resp.Error.Details = []errorDetails{
							{
								ErrorCode: errorCodeUnspecified,
							},
						}
					})

					It("should return general error", func() {
						expectedErr = fmt.Errorf("unsuccessful sendResponse with status code: %d: %s",
							statusCode, string(respBody))

						err := handleResponse(statusCode, respBody)
						Ω(err).Should(Equal(expectedErr))
					})
				})
			})
		})
	})
})
