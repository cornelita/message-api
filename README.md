# CRUD - Message API

## Deployment

[https://law-message-api.herokuapp.com](https://law-message-api.herokuapp.com)

## Endpoint

1. **GET** `{{base_url}}/message/list`

    ![image](./example/get_all_message.png)

2. **GET** `{{base_url}}/message/:id`

    * Path Variables:
        | Key | Type   | Required? |
        | --- | ------ | --------- |
        | id  | Number | yes       |

    ![image](./example/get_message_by_id.png)

3. **POST** `{{base_url}}/message/create`

    * Content-type:
        * multipart/form-data
    * Request Body:
        | Key      | Type   | Required? |
        | -------- | ------ | --------- |
        | name     | String | yes       |
        | password | String | yes       |
        | data     | String | yes       |

    ![image](./example/create_message.png)

4. **PUT** `{{base_url}}/message/:id`

    * Path Variables:
        | Key | Type   | Required? |
        | --- | ------ | --------- |
        | id  | Number | yes       |
    * Content-type:
        * multipart/form-data
    * Request Body:
        | Key      | Type   | Required? |
        | -------- | ------ | --------- |
        | name     | String | yes       |
        | password | String | yes       |
        | data     | String | yes       |

    ![image](./example/update_message.png)

5. **DELETE** `{{base_url}}/message/:id`

    * Path Variables:
        | Key | Type   | Required? |
        | --- | ------ | --------- |
        | id  | Number | yes       |
    * Request Headers:
        | Key           | Value                | Required? |
        | ------------- | -------------------- | --------- |
        | Authorization | <message's password> | yes       |

    ![image](./example/delete_message.png)

## How to Run (Local Windows)

Download [simple-crud.exe](./simple-crud.exe) and run it 🎉

## Postman Collection

[Message API's Postman Collection](./example/Simple%20CRUD.postman_collection.json)
