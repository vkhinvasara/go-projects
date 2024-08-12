package main
import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/gin-gonic/gin"
    "net/http"
)


var db *dynamodb.DynamoDB

func init(){
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)
	if err != nil {
		panic(err)
	}
	db = dynamodb.New(sess)
}

func main(){
	r := gin.Default()
	r.POST("/create", createItem)
	r.GET("/read/:id", readItem)
	r.PUT("/update", updateItem)
	r.Run(":8080")
}
type Item struct{
	ID string`json:"id"`
	Name string `json:"name"`
}
func createItem(c *gin.Context){
	var Item Item
	if err := c.ShouldBindJSON(&Item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(Item.ID),
			},
			"name": {
				S: aws.String(Item.Name),
			},
		},
		TableName: aws.String("users"),
	}
	_, err := db.PutItem(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item created"})
	
}

func readItem(c *gin.Context){
	id := c.Param("id")
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String("users"),
	}
	result, err := db.GetItem(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result.Item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	item := Item{
		ID: *result.Item["id"].S,
		Name: *result.Item["name"].S,
	}
	c.JSON(http.StatusOK, item)
}

func updateItem(c *gin.Context){
	var Item Item
	if err := c.ShouldBindJSON(&Item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(Item.Name),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(Item.ID),
			},
		},
		TableName: aws.String("users"),
		UpdateExpression: aws.String("set name = :n"),
	}
	_, err := db.UpdateItem(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item updated"})
}
