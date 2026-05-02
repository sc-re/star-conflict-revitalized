# Network Message Format for Shard/Loadbalancer Communication

All messages follow the same length encoded format and most are big endian coded.

## Message Header

```mermaid
---
title: CMD Header
---
packet
0-31: "Length (Bytes)"
32-47: "Sequence"
48-63: "Return Sequence"
64-79: "Command Type"
80-95: "Murmur Hash"
```

| Field           | Description                                                                         |
|-----------------|-------------------------------------------------------------------------------------|
| Length          | Message Length in Bytes                                                             |
| Sequence        | Incremental Sequence Number of the Message Sender                                   |
| Return Sequence | Set to the Sequence of the Request Package (if it exists)                           |
| Command Type    | Command Type of the Message                                                         |
| Murmur Hash     | Murmur v2 Hash over the entire Message (header seems to be hashed as little endian) |

## Command Types

- `SCMD_*` -> Server to Client Message
- `CCMD_*` -> Client to Server Message
- `CSCMD_*` -> Client to Server to Client Roundtrip Message

| Idx | Name                             | Idx | Name                              |
|----:|----------------------------------|----:|-----------------------------------|
|   0 | `SCMD_ASSIGNED_SHARD`            |  20 | `SCMD_QUEST_NOTIFICATION`         |
|   1 | `SCMD_LB_QUEUE_INFO`             |  21 | `SCMD_LEAGUE_NOTIFICATION`        |
|   2 | `SCMD_LB_CVARS`                  |  22 | `SCMD_VESSEL_NOTIFICATION`        |
|   3 | `SCMD_AUTH_REQ`                  |  23 | `SCMD_LOBBY_NOTIFICATION`         |
|   4 | `CCMD_AUTH_REQUEST`              |  24 | `SCMD_KEEP_ALIVE`                 |
|   5 | `SCMD_AUTH_ACK`                  |  25 | `SCMD_BAN_INFO`                   |
|   6 | `SCMD_STEAM_NOT_ATTACHED`        |  26 | `SCMD_WELCOME_MSG`                |
|   7 | `SCMD_ARC_NOT_ATTACHED`          |  27 | `SCMD_DOCK_SPACE_STATION`         |
|   8 | `CCMD_STORE`                     |  28 | `SCMD_FREE_SPACE_DEBRIEFING`      |
|   9 | `SCMD_STORE`                     |  29 | `SCMD_NEW_MOTD`                   |
|  10 | `SCMD_STORE_SPOILED`             |  30 | `SCMD_TOURNAMENT_TEAMS_INFO`      |
|  11 | `SCMD_CONNECT_DEDICATED_SERVER`  |  31 | `SCMD_BRAWL_SCHEDULE`             |
|  12 | `SCMD_GAME_ENDED`                |  32 | `SCMD_REWARD_SCHEDULE`            |
|  13 | `CSCMD_ASYNC_REQ`                |  33 | `SCMD_PVE_SCHEDULE`               |
|  14 | `SCMD_NOTIFICATION`              |  34 | `SCMD_LEAGUE_FORBIDDEN_EQUIPMENT` |
|  15 | `SCMD_SQUAD_NOTIFICATION`        |  35 | `SCMD_BATTLE_PASS_ACTIVATION`     |
|  16 | `SCMD_SOCIAL_NOTIFICATION`       |  36 | `SCMD_ZONES_WITH_DISABLED_QUESTS` |
|  17 | `SCMD_TEACH_NOTIFICATION`        |  37 | `SCMD_ADVENTURE_NOTIFICATION`     |
|  18 | `SCMD_CLAN_NOTIFICATION`         |  38 | `SCMD_REPLACE_CHAT_MSG`           |
|  19 | `SCMD_USER_PROFILE_NOTIFICATION` |     |                                   |

