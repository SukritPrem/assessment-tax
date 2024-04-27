package calculateTax

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"fmt"
  "io/ioutil"
  "strings"
  "encoding/json"
  "github.com/go-playground/validator/v10"
  // "errors"
)

type Request_amount struct {
  Amount *float64 `json:"amount" validate:"required,validateValueFloat"`
}

type Response_amount_personalDeduction struct {
  Amount float64 `json:"personalDeduction"`
}

type Response_amount_kReceipt struct {
  Amount float64 `json:"k-receipt"`
}

func (h *Handler) DeductionsPersonal(c echo.Context) error {
  body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
  }
	defer c.Request().Body.Close()

  amount := new(Request_amount)
  amount, err = validateReqAdmin(body)
  if err != nil {
    return c.JSON(http.StatusBadRequest, err.Error())
  }

  if (*amount.Amount <= 10000 || *amount.Amount > 100000){
    return c.JSON(http.StatusBadRequest, "Amount is not in range")
  }
  _, err = h.store.UpdateAmountByTaxType("personalDeduction",*amount.Amount)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, "Not found")
  }
  r := &Response_amount_personalDeduction{
    Amount: *amount.Amount,
  }
  return c.JSON(http.StatusOK, r)
}

func (h *Handler) DeductionsKReceipt(c echo.Context) error {
  body, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
      return err
  }
  defer c.Request().Body.Close()

  amount := new(Request_amount)
  amount, err = validateReqAdmin(body)
  if err != nil {
    return c.JSON(http.StatusBadRequest, err.Error())
  }

  if (*amount.Amount <= 0 || *amount.Amount > 100000){
    return c.JSON(http.StatusBadRequest, "Amount is not in range")
  }
  _, err = h.store.UpdateAmountByTaxType("k-receipt",*amount.Amount)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, "Not found")
  }
  r := &Response_amount_kReceipt{
    Amount: *amount.Amount,
  }
  return c.JSON(http.StatusOK, r)
}



func validateReqAdmin(jsonBytes []byte) (*Request_amount, error){
  amount := new(Request_amount)
  err := CheckDuplicateKey(jsonBytes)
  if err != nil {
    return amount, err
  }
  amount, err = UnmarshalAndValidateReqAdmin(jsonBytes)
  if err != nil {
    return amount, err
  }
  return amount, nil
}

func UnmarshalAndValidateReqAdmin(jsonBytes []byte) (*Request_amount, error) {
  amount := Request_amount{}
  err := validateKeyReqAdmin(jsonBytes)
  if err != nil {
    return &amount, err
  }
  err = json.Unmarshal(jsonBytes, &amount)
  if err != nil {
    return &amount, err
  } 

  validate := validator.New()
  validate.RegisterValidation("validateValueFloat", validateValueFloat)
  err = validate.Struct(amount)
  if err != nil {
    errors := err.(validator.ValidationErrors)
    allErrors  := "Error:"
    for _, e := range errors {
      allErrors = allErrors + e.Field() + " " + e.Tag()
      return &amount, fmt.Errorf(allErrors)
    }
  }
  return &amount, err
}

func CheckDuplicateKey(jsonBytes []byte) error {
  err := check(json.NewDecoder(strings.NewReader(string(jsonBytes))), nil, dupErr)
  if err != nil {
    return fmt.Errorf("Duplicate key found")
  }
  return nil
}












// First Solution I think can use param to check taxType but
// when I read a subject again I think it's not work because
// In subject want /admin/deductions/personal and /admin/deductions/k-receipt
// not /admin/deductions/:taxType
// func (h *Handler) Deductions(c echo.Context) error {
//   taxType := c.Param("taxType")
//   fmt.Println(taxType)
//   body, err := ioutil.ReadAll(c.Request().Body)
// 		if err != nil {
// 			return err
//   }
// 	defer c.Request().Body.Close()
//   var jsonBytes []byte
//   jsonBytes = body
//   err = check(json.NewDecoder(strings.NewReader(string(jsonBytes))), nil, dupErr)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	}
//   err = validateKeyReqAdmin(body)
//   if err != nil {
//     return c.JSON(http.StatusBadRequest, err.Error())
//   }
//   amount := new(Request_amount)
//   validate := validator.New()
//   validate.RegisterValidation("maxfloat", validateValueFloat)
// 	err = validate.Struct(amount)
//   if err != nil {
// 		errors := err.(validator.ValidationErrors)
// 		allErrors := errors.Error()
// 		for _, e := range errors {
// 			allErrors = allErrors + e.Field() + " " + e.Tag() + "\n"
// 			fmt.Println(e.Field(), e.Tag())
// 		}
// 	}
//   if(taxType == "personal"){
//     err = json.Unmarshal(jsonBytes, &amount)
//     if err != nil {
//       return c.JSON(http.StatusBadRequest, "Invalid JSON data")
//     }
//     if (*amount.Amount <= 10000 || *amount.Amount > 100000){
//       return c.JSON(http.StatusBadRequest, "Amount is not in range")
//     }
//     _, err = h.store.UpdateAmountByTaxType("personalDeduction",*amount.Amount)
//     if(err != nil){
//       return c.JSON(http.StatusBadRequest, "Not found")
//     }
//     taxType = "personalDeduction"
//   } else if(taxType == "k-receipt"){
//     err = json.Unmarshal(jsonBytes, &amount)
//     if err != nil {
//       return c.JSON(http.StatusBadRequest, "Invalid JSON data")
//     }
//     if (*amount.Amount <= 0 || *amount.Amount > 100000){
//       return c.JSON(http.StatusBadRequest, "Amount is not in range")
//     }
//     _, err = h.store.UpdateAmountByTaxType("k-receipt",*amount.Amount)
//     if(err != nil){
//       return c.JSON(http.StatusBadRequest, "Not found")
//     }
//     taxType = "kReceipt"
//   } else {  
//     return c.JSON(http.StatusBadRequest, "Invalid tax type")
//   }
// 	data := map[string]interface{}{
//     taxType: *amount.Amount,
//   }
//   jsonData, err := json.Marshal(data)
//   if err != nil {
//     return c.JSON(http.StatusBadRequest, "Invalid JSON data")
//   }
//   return c.JSON(http.StatusOK, jsonData)
// }