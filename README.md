# PoPlan - Planning Poker

## Register

```http
POST http://localhost/register
Content-Type: application/json

{
  "name": "Your Name"
}
```

Response:

```json
{
  "user": {
    "id": "6fe8b74d-b8ce-45de-a3e3-119f752c1e0e",
    "name": "Your Name"
  },
  "token": "6fe8b74d-b8ce-45de-a3e3-119f752c1e0e"
}
```
`token` will be used for authorization.


## Create room

```http
POST http://localhost/rooms
Content-Type: application/json
Authorization: 6fe8b74d-b8ce-45de-a3e3-119f752c1e0e

{
  "name": "Room 1",
  "voteTemplate": 1
}
```

<details><summary>Room response example</summary>

```json
{
  "id": "t0m9f7ar",
  "name": "Room 1",
  "status": 1,
  "seats": [
    {
      "user": {
        "id": "6fe8b74d-b8ce-45de-a3e3-119f752c1e0e",
        "name": "Your Name"
      },
      "vote": {
        "value": 0,
        "type": "unknown"
      },
      "voted": false,
      "owner": true,
      "active": false
    }
  ],
  "template": {
    "title": "Modified fibonacci",
    "votes": [
      {
        "value": 0,
        "type": "unknown"
      },
      {
        "value": 0,
        "type": "value"
      },
      {
        "value": 0.5,
        "type": "value"
      },
      {
        "value": 1,
        "type": "value"
      }
    ]
  }
}
```
</details>

## Get templates to create room

```http
GET http://localhost/rooms/templates
Content-Type: application/json
Authorization: 6fe8b74d-b8ce-45de-a3e3-119f752c1e0e
```
<details><summary>Response example</summary>

```json
{
  "templates": [
    {
      "title": "Fibonacci",
      "votes": [
        {
          "value": 0,
          "type": "unknown"
        },
        {
          "value": 0,
          "type": "value"
        },
        {
          "value": 1,
          "type": "value"
        },
        {
          "value": 0,
          "type": "infinity"
        }
      ]
    },
    {
      "title": "Modified fibonacci",
      "votes": [
        {
          "value": 0,
          "type": "unknown"
        },
        {
          "value": 0,
          "type": "break"
        },
        {
          "value": 0,
          "type": "value"
        },
        {
          "value": 0.5,
          "type": "value"
        },
        {
          "value": 1,
          "type": "value"
        }
      ]
    }
  ]
}
```
</details>

Index should be passed to `voteTemplate` field when creating a room.

## Get room

Join room and get info about it.

```http
GET http://localhost/rooms/t0m9f7ar
Content-Type: application/json
Authorization: 6fe8b74d-b8ce-45de-a3e3-119f752c1e0e
```
Response will contain room.

## Update room

Only owner can update room.

```http
PATCH http://localhost/rooms/t0m9f7ar
Content-Type: application/json
Authorization: 6fe8b74d-b8ce-45de-a3e3-119f752c1e0e

{
  "name": "New room name"
}
```

Response will contain room.

## Leave room

```http
POST http://localhost/rooms/t0m9f7ar/leave
Content-Type: application/json
Authorization: 6fe8b74d-b8ce-45de-a3e3-119f752c1e0e
```

Response:
```json
{
  "message": "Room left"
}
```

## Vote

```http
POST http://localhost/rooms/t0m9f7ar/vote
Content-Type: application/json
Authorization: 6fe8b74d-b8ce-45de-a3e3-119f752c1e0e

{
  "value": 5
}
```

`value` should contain index from vote template.

Response:
```json
{
  "message": "Voted"
}
```

## End vote

```http
POST http://localhost/rooms/t0m9f7ar/end
Content-Type: application/json
Authorization: 6fe8b74d-b8ce-45de-a3e3-119f752c1e0e
```

Only owner can do this.

Response will contain room.

## Reset room

```http
POST http://localhost/rooms/t0m9f7ar/reset
Content-Type: application/json
Authorization: 6fe8b74d-b8ce-45de-a3e3-119f752c1e0e
```

Only owner can do this.

Response will contain room.

## Subscribe

Example shown in [example file](example/index.html)

On room status change, voting, joining and leaving message will be received with room object.

## Room object

<details><summary>Room response example</summary>

```json
{
  "id": "t0m9f7ar",
  "name": "Room 1",
  "status": 1,
  "seats": [
    {
      "user": {
        "id": "6fe8b74d-b8ce-45de-a3e3-119f752c1e0e",
        "name": "Your Name"
      },
      "vote": {
        "value": 0,
        "type": "unknown"
      },
      "voted": false,
      "owner": true,
      "active": false
    }
  ],
  "template": {
    "title": "Modified fibonacci",
    "votes": [
      {
        "value": 0,
        "type": "unknown"
      },
      {
        "value": 0,
        "type": "value"
      },
      {
        "value": 0.5,
        "type": "value"
      },
      {
        "value": 1,
        "type": "value"
      }
    ]
  }
}
```
</details>

  * `id` - id, used in url
  * `name` - name to show
  * `status` - room status, 1 - voting, 2 - voted, triggered by **end vote** and **reset room**
  * `seats` - users in the room
    * `user` - info about user in seat
    * `vote` - user's vote
      * `value` - integer value
      * `type` - vote type, available types:
        * `value` - numeric value
        * `unknown` - vote is unknown, "?"
        * `break` - need a break, "coffee"
        * `infinity` - infinity
    * `voted` - user has voted
    * `owner`- user is room owner
    * `active` - user is active, subscribed and will get all the info
  * `template` - vote template with list of available options
