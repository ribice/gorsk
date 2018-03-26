package mw

// Error handling middleware for Gin
// I prefer using a more generic approach (apperr.Response) since it makes it easier to migrate from Gin, is easier to notice for other users and not all requests have to go through it
// // CatchError catches all errors and creates JSON response from them
// func CatchError() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		defer func() {
// 			if len(c.Errors) > 0 {
// 				err := c.Errors[0].Err
// 				switch err.(type) {
// 				case *apperr.APPError:
// 					e := err.(*apperr.APPError)
// 					c.AbortWithStatusJSON(e.Status, e)
// 					return
// 				case validator.ValidationErrors:
// 					var errMsg []string
// 					e := err.(validator.ValidationErrors)
// 					for _, v := range e {
// 						errMsg = append(errMsg, fmt.Sprintf("%s: condition %s not satisfied", v.Name, v.ActualTag))
// 					}
// 					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errMsg})
// 				default:
// 					c.AbortWithStatus(http.StatusInternalServerError)
// 				}
// 			}
// 		}()
// 		c.Next()
// 	}
// }
