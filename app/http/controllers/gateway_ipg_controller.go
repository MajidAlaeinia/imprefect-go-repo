package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/vandario/govms-ipg/app/http/validations"
	"github.com/vandario/govms-ipg/database"
	"net/http"
	"strconv"
	"time"
)

type RequestFields struct { //TODO: more strict validation?
	ApiKey       string `json:"api_key" binding:"required"`
	Amount       int    `json:"amount" binding:"required,min=1000,numeric"`
	CallbackUrl  string `json:"callback_url" binding:"required,url"`
	MobileNumber string `json:"mobile_number"`
	FactorNumber string `json:"factorNumber"`
	Description  string `json:"description"`
}

var requestFields RequestFields

type GatewayIpg struct {
	gorm.Model
	UserId     int
	BusinessId int
	//Url string `json:"url"`
	Urls   string
	Ips    string
	ApiKey string
	Status int
	Review bool
	//Ip []byte `json:"ip"`
	WageType bool
	//IsVip bool `json:"is_vip"`
	//HasPayedForAdditionalBusiness bool `json:"has_payed_for_additional_business"`
	MerchantId       string
	MerchantPassword string
	//WagePercent string `json:"wage_percent"`
	//WageMin string `json:"wage_min"`
	//WageMax string `json:"wage_max"`
}

var gatewayIpg GatewayIpg

type GatewayIpgLog struct {
	ID           int    `gorm:"primary_key"`
	Token        string `gorm:"unique;not null"`
	ApiKey       string
	CallbackUrl  string
	Amount       string // as in current database.
	Ip           string
	FactorNumber string `gorm:"column:factorNumber"`
	Description  string `gorm:"text"` //TODO: Right?
	MobileNumber string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var gatewayIpgLog GatewayIpgLog

type GatewayTransaction struct {
	ID    int    `gorm:"primary_key"`
	Port  string //TODO: Enum?
	Price int    //TODO: Decimal?
	//RefId                string
	//TrackingCode         string
	//CardNumber           string
	//Cid                  string
	Status string //TODO: Enum?
	//RevisedTransactionId int
	//Confirmed            bool
	Ip string
	//PaymentDate          time.Time
	//PaymentNumber        string
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt            time.Time
	//RequestFormsId       int
	MobileNumber string
	UserId       int
	IpgId        uint
	BusinessId   int
	Channel      string
	CallbackUrl  string
	FactorNumber string `gorm:"column:factorNumber"`
	Description  string //TODO: text?
	ApiToken     string
	//Response             bool
	Wage string
	//IbanId               string // as in the existing database.
	//SourceNumber         string
	//SettlementMessage    string
	//SettlementState      string
	//SettlementType       string
	//InquirySequence      string
	//InquiryDate          string
	//InquiryTime          string
	//Wallet               string // as in the existing database.
	//Address              string // as in the existing database.
	//Email                string
	//Name                 string
	//Phone                string
	//MerchantId           string
	//MerchantPassword     string
	//VandarBalance        string
	//VandarAccountId      int
	//UpdatedBy            int
	//DeletedBy            int
}

//var gatewayTransaction GatewayTransaction

func GatewayIpgSendHandler(context *gin.Context) {
	// Connect to the database
	vandarDatabase := database.VandarDatabase()

	// Bind the requested json
	bindJsonError := context.BindJSON(&requestFields)

	// Validation Errors
	if bindJsonError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": bindJsonError.Error(),
		})
		return
	}

	// Find 'GatewayIpg' according to the 'api_key' existence in 'GatewayIpg'
	findGatewayIpgError := vandarDatabase.Model(&gatewayIpg).Where("api_key = ?", requestFields.ApiKey).Find(&gatewayIpg).Error
	if findGatewayIpgError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": findGatewayIpgError.Error(),
		})
		return
	}

	// Custom Validations

	// Callback Url Validation
	isCallbackUrlValid := validations.CallbackUrlValidation(gatewayIpg.Urls, requestFields.CallbackUrl)

	if !isCallbackUrlValid {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Callback Url is invalid.", //TODO: A better message.
		})
		return
	}

	// Create the validated 'GatewayIpgLog'
	token := uuid.New().String()
	gatewayIpgLog := GatewayIpgLog{
		Token:        token,
		ApiKey:       requestFields.ApiKey,
		CallbackUrl:  requestFields.CallbackUrl,
		Amount:       strconv.Itoa(requestFields.Amount),
		Ip:           context.ClientIP(),
		FactorNumber: requestFields.FactorNumber,
		Description:  requestFields.Description,
		MobileNumber: requestFields.MobileNumber,
	}
	vandarDatabase.Create(&gatewayIpgLog)

	context.JSON(http.StatusOK, gin.H{
		"status": 1,
		"token":  token,
	})
	return
}

func GatewayIpgRequestHandler(context *gin.Context) {
	// Connect to the database
	vandarDatabase := database.VandarDatabase()

	token := context.Param("token")
	findGatewayIpgLogError := vandarDatabase.Model(&gatewayIpgLog).Where("token = ?", token).Find(&gatewayIpgLog).Error
	if findGatewayIpgLogError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": findGatewayIpgLogError.Error(),
		})
		return
	}

	findGatewayIpgError := vandarDatabase.Model(&gatewayIpg).Where("api_key = ?", &gatewayIpgLog.ApiKey).Find(&gatewayIpg).Error
	if findGatewayIpgError != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": findGatewayIpgError.Error(),
		})
		return
	}
	wage := CalculateWage(gatewayIpgLog.Amount)
	gatewayTransactionId, _ := strconv.Atoi(GetTimeId())
	amount, _ := strconv.Atoi(gatewayIpgLog.Amount)

	gatewayTransaction := GatewayTransaction{
		ID:           gatewayTransactionId,
		Port:         "SAMAN",
		Price:        amount,
		Status:       "INIT",
		Ip:           context.ClientIP(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		MobileNumber: gatewayIpgLog.MobileNumber,
		UserId:       gatewayIpg.UserId,
		IpgId:        gatewayIpg.ID,
		BusinessId:   gatewayIpg.BusinessId,
		Channel:      "IPG",
		CallbackUrl:  gatewayIpgLog.CallbackUrl,
		FactorNumber: gatewayIpgLog.FactorNumber,
		Description:  gatewayIpgLog.Description,
		ApiToken:     gatewayIpgLog.Token,
		Wage:         wage,
		//MerchantId: ,
		//MerchantPassword:
	}
	vandarDatabase.Create(&gatewayTransaction)
	vandarDatabase.Model(&gatewayIpgLog).Delete(&gatewayIpgLog)

	//TODO: Redirect with wanted variables.
	context.Request.URL.Path = "/index"
	apiRoutes := gin.Default()
	apiRoutes.HandleContext(context)

	context.Redirect(http.StatusFound, "http://www.google.com/")

	//fmt.Println("DONE")
	//os.Exit(1)
}

func CalculateWage(amount string) string {
	var wagePercentage = 0.01 //TODO: use env.
	amountFloat, _ := strconv.ParseFloat(amount, 8)
	wage := wagePercentage * amountFloat
	if wage >= 30000 { //TODO: Use env.
		wage = 30000 //TODO: Use env.
	}
	return strconv.FormatFloat(wage, 'f', 6, 64)
}

func GetTimeId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)[0:12]
}
