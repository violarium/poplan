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
  "id": "831a6wkf",
  "name": "Room 1",
  "status": 1,
  "owner": true,
  "seats": [
    {
      "user": {
        "id": "6a4be959-d3f7-469a-8b8a-820fd5cd2081",
        "name": "Your Name"
      },
      "vote": {
        "value": 0,
        "type": "unknown"
      },
      "voted": false,
      "voteOpened": true,
      "owner": true,
      "active": false
    }
  ],
  "templateTitle": "Modified fibonacci",
  "voteCards": [
    {
      "vote": {
        "value": 0,
        "type": "unknown"
      },
      "active": false
    },
    {
      "vote": {
        "value": 0,
        "type": "value"
      },
      "active": false
    },
    {
      "vote": {
        "value": 0.5,
        "type": "value"
      },
      "active": false
    },
    {
      "vote": {
        "value": 1,
        "type": "value"
      },
      "active": false
    },
    {
      "vote": {
        "value": 2,
        "type": "value"
      },
      "active": false
    },
    {
      "vote": {
        "value": 0,
        "type": "infinity"
      },
      "active": false
    }
  ]
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
  "id": "831a6wkf",
  "name": "Room 1",
  "status": 1,
  "owner": true,
  "seats": [
    {
      "user": {
        "id": "6a4be959-d3f7-469a-8b8a-820fd5cd2081",
        "name": "Your Name"
      },
      "vote": {
        "value": 0,
        "type": "unknown"
      },
      "voted": false,
      "voteOpened": true,
      "owner": true,
      "active": false
    }
  ],
  "templateTitle": "Modified fibonacci",
  "voteCards": [
    {
      "vote": {
        "value": 0,
        "type": "unknown"
      },
      "active": false
    },
    {
      "vote": {
        "value": 0,
        "type": "value"
      },
      "active": false
    },
    {
      "vote": {
        "value": 0.5,
        "type": "value"
      },
      "active": false
    },
    {
      "vote": {
        "value": 1,
        "type": "value"
      },
      "active": false
    },
    {
      "vote": {
        "value": 2,
        "type": "value"
      },
      "active": false
    },
    {
      "vote": {
        "value": 0,
        "type": "infinity"
      },
      "active": false
    }
  ]
}
```
</details>

  * `id` - id, used in url
  * `name` - name to show
  * `status` - room status, 1 - voting, 2 - voted, triggered by **end vote** and **reset room**
  * `owner` - is current user (you) is owner of this room
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
    * `voteOpened` - is card opened or not?
    * `owner`- user is room owner
    * `active` - user is active, subscribed and will get all the info
  * `templateTitle` - title of the room template
  * `voteCards` - available vote cards
    * `vote` - vote, same as within template
    * `active` - card is active (selected) for current user (you)
