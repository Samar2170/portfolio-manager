Feature: login
    In order to use this service you need to login

    Scenario: Login using username and password
        Given i signup
        Then i login 
        