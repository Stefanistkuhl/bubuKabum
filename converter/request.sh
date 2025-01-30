curl -X POST http://localhost:6999/api/emote \
  -H "Content-Type: application/json" \
  -d '{
    "emotes": [
      {
        "link": "https://7tv.app/emotes/01FW23V83G000EFCFJ99X9MAQZ",
        "is_2_frame_gif": true,
        "desired_name": "trash"
      },
      {
        "link": "https://7tv.app/emotes/01F6MKTFTG0009C9ZSNZTFV2ZF",
        "is_2_frame_gif": false,
        "desired_name": ""
      }
    ]
  }'

