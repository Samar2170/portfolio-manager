### Design

#### Data Models  (Trade / Holding) 
* Trade is a flow variable and Holding is a static variable. 

1. Security Types - Attributes 
    * Stocks - Units, Price/Unit, Total, CV
    * Bonds - Units, Price/Unit, Total, IPFreq, IPDate, IPRate, MtDate, IsCallable?
        * Monthly, Annually
        * Zero coupon
        * Callable
        * Partially Maturing

    * Leases 
    * FD - Total, IPRate, IPFreq, IPDate, MtDate
    * Crypto - Units, Price/Unit, Total, CV
    * MF - Units, Price/Unit, Total, CV
    * Unlisted Govt Bonds - Units, Price/Unit, Total, CV


### API Endpoints

    ----------|User|-----------------
    |                               |
    |                               |
____________________        ___________________
| Trade Reg Service|        | Securities Serv   |






### Ancillary Services
* Bank/Cash Mgmt Service
* Security Names Service / Company Price Update Service
* Messenging Service (Signup using a token, that is to be given by the api and registered by messaging on telegram chat.)



##
* DO less at the api level , more over controller. WIll help in reusing code in bulk upload.


### Functionality
## FD 
* Calculate next ip_date and update it. (Cronj, securities)
* Notify When interest is due.  (Cronj, portfolio)
* Mark Instruments as expired once Maturity date is up. (Cronj, portfolio)
* Calculate Accrued Interest.   (Cronj, portfolio)



#### Telegram Integration
##### ADR
* How to implement ? 
    1. Telegram parses message [View,Register]
        1. Level 2 Parser parses formname inside [`[RegisterDematAccount]`]
        2. Get params keys and values, Create a JSON , send it to api endpoint. 
        3. Write Auth Middleware for http client for this. 

* Auth, API gives out a 6 Digit OTP, which is then used in telegram to authenticate.

* Endpoints
    * Authenticate using OTP.
    * Corresponding endpoint for each API endpoint. 

