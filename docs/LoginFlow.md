# LoginFlow

## Authentication (Username + Password)

On a higher level this is just a simple DiffieHellman KEX to derive an AES key to encrypt the password and send it together with username and some Client Information to the Server.
The server responds with either an `SCMD_AUTH_ACK` or with a Disconnect Message.
* The flow is different for Accounts with 2FA enabled or Steam Login


## Post Authentication

After Authentication we should trigger a `SCMD_USER_PROFILE_NOTIFICATION` and inform all Clients that either follow or have the user befriended about the changed online state.

The client is now in a limbo state and will only request a few `AC_` types and otherwise wait for the server to send a `AC_VESSEL_STRIP_IMPROPER_BATTLE`, which is unusual.
Most `CSCMD_ASYNC_REQ` Messages are requested by the user, with only this exception and another:

- `AC_VESSEL_STRIP_IMPROPER_BATTLE`: This message sent after the sever checked that all stored ships only have valid equipment and lists all ships that had to be stripped. It also contains sends the account exp/clearance score.
- `AC_UPDATE_YUP_PURCHASES`: Contains a list of DLCs purchases by the user, only sent if this list wouldn't be empty.


