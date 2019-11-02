*API Description*
----

**Server's Health Check**
----
  Return minimal JSON information about server.

* **URL**

  /api/v1/server/status

* **Method:**
  
  `GET`
  
*  **URL Params**

   None

* **Data Params**

  None

* **Success Response:**

  * **Code:** 200  
    **Content:** `{ "is_alive": true, "build_time": "20191030.191528", "version": "8c53876f7dd2f8fece47b01eb1c219ace4bdf183" }`
 
* **Error Response:**

  None

* **Sample Call:**

  ```sh
  curl -X GET http://localhost:8080/api/v1/server/status
  ``` 

**Create Account**
----
  Create new user's account.

* **URL**

  /api/v1/accounts

* **Method:**
  
  `POST`
  
*  **URL Params**

   None

* **Data Params**

  Descriptions of the new account.
  
  ```json
    {
        "uid": "toshik1978",
        "currency": "USD",
        "balance": 100
    }
  ```

* **Success Response:**
  
  Account created.

  * **Code:** 200 <br />
    **Content:** `{ "uid": "toshik1978", "currency": "USD", "balance": 100, "created_at": "2019-11-02T20:29:18.76046542Z" }`
 
* **Error Response:**

  * **Code:** 400 BAD REQUEST  
    **Content:** `failed to create account: uid is empty`

  OR

  * **Code:** 500 INTERNAL SERVER ERROR  
    **Content:** `failed to create account: database failure`

* **Sample Call:**

  ```sh
    curl -X POST \
      http://localhost:8080/api/v1/accounts \
      -H 'Content-Type: application/json' \
      -d '{
            "uid": "toshik1978",
            "currency": "USD",
            "balance": 100
          }'
  ```

**Get All Accounts**
----
  Get all user accounts.

* **URL**

  /api/v1/accounts

* **Method:**
  
  `GET`
  
*  **URL Params**

   None

* **Data Params**

  None

* **Success Response:**
  
  List of all accounts.

  * **Code:** 200 <br />
    **Content:** `[{ "uid": "toshik1978", "currency": "USD", "balance": 100, "created_at": "2019-11-02T20:29:18.760465Z" }]`
 
* **Error Response:**

  * **Code:** 500 INTERNAL SERVER ERROR  
    **Content:** `failed to get accounts: database failure`

* **Sample Call:**

  ```sh
    curl -X GET http://localhost:8080/api/v1/accounts
  ```

**Create Payment**
----
  Create new payment from one account to another.

* **URL**

  /api/v1/accounts/toshik1978/payments

* **Method:**
  
  `POST`
  
*  **URL Params**

   None

* **Data Params**

  Descriptions of the new account.
  
  ```json
    {
        "recipient": "toshik1979",
        "amount": 100
    }
  ```

* **Success Response:**
  
  Payment created.

  * **Code:** 200 <br />
    **Content:** `{ "account": "toshik1978", "to_account": "toshik1979", "direction": "outgoing", "amount": 100, "created_at": "2019-11-02T20:30:52.374818264Z" }`
 
* **Error Response:**

  * **Code:** 400 BAD REQUEST  
    **Content:** `failed to create payment: uid is empty`

  OR

  * **Code:** 500 INTERNAL SERVER ERROR  
    **Content:** `failed to create payment: database failure`

* **Sample Call:**

  ```sh
    curl -X POST \
      http://localhost:8080/api/v1/accounts/toshik1978/payments \
      -H 'Content-Type: application/json' \
      -d '{
            "recipient": "toshik1979",
            "amount": 100
          }'
  ```

**Get All Payments**
----
  Get all payments.

* **URL**

  /api/v1/accounts/payments

* **Method:**
  
  `GET`
  
*  **URL Params**

   None

* **Data Params**

   None

* **Success Response:**
  
  List of all payments.

  * **Code:** 200 <br />
    **Content:** `[{ "account": "toshik1978", "to_account": "toshik1979", "direction": "outgoing", "amount": 100, "created_at": "2019-11-02T20:30:52.374818Z" },
                      { "account": "toshik1979", "from_account": "toshik1978", "direction": "incoming", "amount": 100, "created_at": "2019-11-02T20:30:52.374818Z" }]`
 
* **Error Response:**

  * **Code:** 500 INTERNAL SERVER ERROR  
    **Content:** `failed to get payments: database failure`

* **Sample Call:**

  ```sh
    curl -X GET http://localhost:8080/api/v1/accounts/payments
  ```
