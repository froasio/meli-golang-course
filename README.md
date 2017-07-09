# meli-golang-course
Final project for Mercado Libre's Golang course

## Description
Develop a REST API written in Golang that calculates the following JSON response given a product category.
```json
{
  "max":10,
  "suggested":5,
  "min":1
}
```
The resource should response on the following route

/categories/$ID/prices

Example:

curl -X GET “http://mydomain.com/categories/MLA3530/prices”
