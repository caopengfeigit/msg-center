/*
 Navicat Premium Data Transfer

 Source Server         : localhost mongodb
 Source Server Type    : MongoDB
 Source Server Version : 40206
 Source Host           : localhost:27017
 Source Schema         : msg_center

 Target Server Type    : MongoDB
 Target Server Version : 40206
 File Encoding         : 65001

 Date: 02/11/2021 17:03:13
*/


// ----------------------------
// Collection structure for callback_request_config
// ----------------------------
db.getCollection("callback_request_config").drop();
db.createCollection("callback_request_config");
db.getCollection("callback_request_config").createIndex({
    "parent_id": NumberInt("1")
}, {
    name: "parent_id_1"
});

// ----------------------------
// Documents of callback_request_config
// ----------------------------
session = db.getMongo().startSession();
session.startTransaction();
db = session.getDatabase("msg_center");
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611cabe18b17902a67626cd3"),
    "parent_id": ObjectId("6088cf0b6fca0ea458770479"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback1",
    "callback_request_type": "GET",
    "callback_request_is_json": false
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611cabe18b17902a67626cd4"),
    "parent_id": ObjectId("6115d3ac2a462317a000a0c3"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback1",
    "callback_request_type": "GET",
    "callback_request_is_json": false
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611cabe18b17902a67626cd5"),
    "parent_id": ObjectId("611caef28b17902a67626cd6"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback1",
    "callback_request_type": "POST",
    "callback_request_is_json": true
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611caff08b17902a67626cd7"),
    "parent_id": ObjectId("60f67ef2c0555f7bce783404"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback1",
    "callback_request_type": "POST",
    "callback_request_is_json": true
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611cb0af8b17902a67626cd9"),
    "parent_id": ObjectId("611cbad88b17902a67626cda"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback2",
    "callback_request_type": "GET",
    "callback_request_is_json": false
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611f52e1fe63421349532874"),
    "parent_id": ObjectId("611f51c1fe63421349532873"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback3",
    "callback_request_type": "GET",
    "callback_request_is_json": true
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611f5890fe6342134953287a"),
    "parent_id": ObjectId("611f5702fe63421349532876"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback4",
    "callback_request_type": "GET",
    "callback_request_is_json": false
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611f5890fe6342134953287b"),
    "parent_id": ObjectId("611f577dfe63421349532877"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback5",
    "callback_request_type": "GET",
    "callback_request_is_json": true
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611f5890fe6342134953287c"),
    "parent_id": ObjectId("611f57a8fe63421349532878"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback6",
    "callback_request_type": "POST",
    "callback_request_is_json": false
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("611f5890fe6342134953287d"),
    "parent_id": ObjectId("611f57f7fe63421349532879"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback7",
    "callback_request_type": "POST",
    "callback_request_is_json": true
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("6131e72dc3d9e37cc124f734"),
    "parent_id": ObjectId("6131d68ac3d9e37cc124f733"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback8",
    "callback_request_is_json": false,
    "callback_request_type": "GET"
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("613346107caeb90e436b9c95"),
    "parent_id": ObjectId("6133354d7caeb90e436b9c93"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback9",
    "callback_request_is_json": true,
    "callback_request_type": "POST"
} ]);
db.getCollection("callback_request_config").insert([ {
    _id: ObjectId("6133463d7caeb90e436b9c96"),
    "parent_id": ObjectId("6133354d7caeb90e436b9c94"),
    "callback_host": "http://testhost.com",
    "callback_path": "/callback/callback9",
    "callback_request_is_json": false,
    "callback_request_type": "POST"
} ]);
session.commitTransaction(); session.endSession();

// ----------------------------
// Collection structure for callback_request_log
// ----------------------------
db.getCollection("callback_request_log").drop();
db.createCollection("callback_request_log");
db.getCollection("callback_request_log").createIndex({
    "created_at": NumberInt("1")
}, {
    name: "created_at_1",
    expireAfterSeconds: NumberInt("604800")
});

// ----------------------------
// Documents of callback_request_log
// ----------------------------
session = db.getMongo().startSession();
session.startTransaction();
db = session.getDatabase("msg_center");
db.getCollection("callback_request_log").insert([ {
    _id: ObjectId("6180a924087f0fe830c2d072"),
    "project_name": "q",
    "event_key": "q",
    "queue_name": "q-1",
    "request_host": "q",
    "request_path": "q",
    "request_type": "GET",
    "callback_request_is_json": false,
    "request_data": "111",
    "request_res": false,
    "request_error": "lookup : no such host",
    "request_status": NumberInt("200"),
    "request_response": "",
    "created_at": ISODate("2021-11-02T02:57:40.25Z"),
    "created_at_str": "2021-11-02 10:57:40"
} ]);
db.getCollection("callback_request_log").insert([ {
    _id: ObjectId("6180a924087f0fe830c2d071"),
    "project_name": "q",
    "event_key": "q",
    "queue_name": "q-1",
    "request_host": "q",
    "request_path": "q",
    "request_type": "GET",
    "callback_request_is_json": false,
    "request_data": "111",
    "request_res": false,
    "request_error": "lookup : no such host",
    "request_status": NumberInt("200"),
    "request_response": "",
    "created_at": ISODate("2021-11-02T02:57:40.249Z"),
    "created_at_str": "2021-11-02 10:57:40"
} ]);
db.getCollection("callback_request_log").insert([ {
    _id: ObjectId("6180a924087f0fe830c2d073"),
    "project_name": "q",
    "event_key": "q",
    "queue_name": "q-1",
    "request_host": "q",
    "request_path": "q",
    "request_type": "GET",
    "callback_request_is_json": false,
    "request_data": "111",
    "request_res": false,
    "request_error": "lookup : no such host",
    "request_status": NumberInt("200"),
    "request_response": "",
    "created_at": ISODate("2021-11-02T02:57:40.25Z"),
    "created_at_str": "2021-11-02 10:57:40"
} ]);
session.commitTransaction(); session.endSession();

// ----------------------------
// Collection structure for project
// ----------------------------
db.getCollection("project").drop();
db.createCollection("project");
db.getCollection("project").createIndex({
    name: NumberInt("1")
}, {
    name: "name_1",
    unique: true
});

// ----------------------------
// Documents of project
// ----------------------------
session = db.getMongo().startSession();
session.startTransaction();
db = session.getDatabase("msg_center");
db.getCollection("project").insert([ {
    _id: ObjectId("60868b316fca0ea458770476"),
    name: "test1",
    "business_line": "业务线1"
} ]);
session.commitTransaction(); session.endSession();

// ----------------------------
// Collection structure for project_event
// ----------------------------
db.getCollection("project_event").drop();
db.createCollection("project_event");
db.getCollection("project_event").createIndex({
    "project_id": NumberInt("1"),
    name: NumberInt("1")
}, {
    name: "project_id_1_name_1",
    unique: true
});
db.getCollection("project_event").createIndex({
    name: NumberInt("1"),
    type: NumberInt("1")
}, {
    name: "name_1_type_1",
    unique: true
});

// ----------------------------
// Documents of project_event
// ----------------------------
session = db.getMongo().startSession();
session.startTransaction();
db = session.getDatabase("msg_center");
db.getCollection("project_event").insert([ {
    _id: ObjectId("6087b7bf6fca0ea458770478"),
    "project_id": ObjectId("60868b316fca0ea458770476"),
    name: "test1",
    type: "Single"
} ]);
db.getCollection("project_event").insert([ {
    _id: ObjectId("60f67e52c0555f7bce783402"),
    "project_id": ObjectId("60868b316fca0ea458770476"),
    name: "test2",
    type: "Single"
} ]);
db.getCollection("project_event").insert([ {
    _id: ObjectId("60f7dea2bd3bc7767f2f4b92"),
    "project_id": ObjectId("60868b316fca0ea458770476"),
    name: "test3",
    type: "PublishSubscribe",
    "exchange_type": "x-delayed-message"
} ]);
db.getCollection("project_event").insert([ {
    _id: ObjectId("611cb0198b17902a67626cd8"),
    "project_id": ObjectId("60868b316fca0ea458770476"),
    name: "test4",
    type: "PublishSubscribe",
    "exchange_type": "fanout"
} ]);
db.getCollection("project_event").insert([ {
    _id: ObjectId("611f513afe63421349532872"),
    "project_id": ObjectId("60868b316fca0ea458770476"),
    name: "test5",
    type: "PublishSubscribe",
    "exchange_type": "direct"
} ]);
db.getCollection("project_event").insert([ {
    _id: ObjectId("611f56b7fe63421349532875"),
    "project_id": ObjectId("60868b316fca0ea458770476"),
    name: "test6",
    type: "PublishSubscribe",
    "exchange_type": "topic"
} ]);
db.getCollection("project_event").insert([ {
    _id: ObjectId("6131d634c3d9e37cc124f732"),
    "project_id": ObjectId("60868b316fca0ea458770476"),
    name: "test7",
    type: "PublishSubscribe",
    "exchange_type": "headers"
} ]);
db.getCollection("project_event").insert([ {
    _id: ObjectId("613324477caeb90e436b9c92"),
    name: "test8",
    "project_id": ObjectId("60868b316fca0ea458770476"),
    type: "WorkQueues"
} ]);
session.commitTransaction(); session.endSession();

// ----------------------------
// Collection structure for publish_subscribe_config
// ----------------------------
db.getCollection("publish_subscribe_config").drop();
db.createCollection("publish_subscribe_config");
db.getCollection("publish_subscribe_config").createIndex({
    "event_id": NumberInt("1")
}, {
    name: "event_id_1"
});
db.getCollection("publish_subscribe_config").createIndex({
    "event_id": NumberInt("1"),
    "exchange_name": NumberInt("1"),
    "routing_key": NumberInt("1"),
    "queue_name": NumberInt("1")
}, {
    name: "event_id_1_exchange_name_1_routing_key_1_queue_name_1",
    unique: true
});

// ----------------------------
// Documents of publish_subscribe_config
// ----------------------------
session = db.getMongo().startSession();
session.startTransaction();
db = session.getDatabase("msg_center");
db.getCollection("publish_subscribe_config").insert([ {
    _id: ObjectId("6115d3ac2a462317a000a0c3"),
    "event_id": ObjectId("60f7dea2bd3bc7767f2f4b92"),
    "exchange_name": "psc_exchange_1",
    "routing_key": "psc_routing_1",
    "queue_name": "psc_queue_1"
} ]);
db.getCollection("publish_subscribe_config").insert([ {
    _id: ObjectId("611caef28b17902a67626cd6"),
    "event_id": ObjectId("611cb0198b17902a67626cd8"),
    "exchange_name": "psc_exchange_2",
    "queue_name": "psc_queue_fanout1",
    "routing_key": null
} ]);
db.getCollection("publish_subscribe_config").insert([ {
    _id: ObjectId("611cbad88b17902a67626cda"),
    "event_id": ObjectId("611cb0198b17902a67626cd8"),
    "exchange_name": "psc_exchange_2",
    "queue_name": "psc_queue_fanout2"
} ]);
db.getCollection("publish_subscribe_config").insert([ {
    _id: ObjectId("611f51c1fe63421349532873"),
    "event_id": ObjectId("611f513afe63421349532872"),
    "exchange_name": "exchange_direct",
    "routing_key": "routing_key_direct",
    "queue_name": "queue_direct"
} ]);
db.getCollection("publish_subscribe_config").insert([ {
    _id: ObjectId("611f5702fe63421349532876"),
    "event_id": ObjectId("611f56b7fe63421349532875"),
    "exchange_name": "exchange_topic",
    "queue_name": "queue_topic_1",
    "routing_key": "topic.*"
} ]);
db.getCollection("publish_subscribe_config").insert([ {
    _id: ObjectId("611f577dfe63421349532877"),
    "event_id": ObjectId("611f56b7fe63421349532875"),
    "exchange_name": "exchange_topic",
    "routing_key": "topic.#",
    "queue_name": "queue_topic_2"
} ]);
db.getCollection("publish_subscribe_config").insert([ {
    _id: ObjectId("611f57a8fe63421349532878"),
    "event_id": ObjectId("611f56b7fe63421349532875"),
    "exchange_name": "exchange_topic",
    "routing_key": "*.topic",
    "queue_name": "queue_topic_3"
} ]);
db.getCollection("publish_subscribe_config").insert([ {
    _id: ObjectId("611f57f7fe63421349532879"),
    "event_id": ObjectId("611f56b7fe63421349532875"),
    "exchange_name": "exchange_topic",
    "routing_key": "#.topic",
    "queue_name": "queue_topic_4"
} ]);
db.getCollection("publish_subscribe_config").insert([ {
    _id: ObjectId("6131d68ac3d9e37cc124f733"),
    "event_id": ObjectId("6131d634c3d9e37cc124f732"),
    "exchange_name": "exchange_headers",
    "queue_name": "queue_headers",
    headers: {
        key1: "123456",
        key2: "abcde"
    },
    "x_match": "any"
} ]);
session.commitTransaction(); session.endSession();

// ----------------------------
// Collection structure for single_config
// ----------------------------
db.getCollection("single_config").drop();
db.createCollection("single_config");
db.getCollection("single_config").createIndex({
    "event_id": NumberInt("1")
}, {
    name: "event_id_1"
});
db.getCollection("single_config").createIndex({
    "event_id": NumberInt("1"),
    "queue_name": NumberInt("1")
}, {
    name: "event_id_1_queue_name_1",
    unique: true
});

// ----------------------------
// Documents of single_config
// ----------------------------
session = db.getMongo().startSession();
session.startTransaction();
db = session.getDatabase("msg_center");
db.getCollection("single_config").insert([ {
    _id: ObjectId("6088cf0b6fca0ea458770479"),
    "event_id": ObjectId("6087b7bf6fca0ea458770478"),
    "queue_name": "hello1"
} ]);
db.getCollection("single_config").insert([ {
    _id: ObjectId("60f67ef2c0555f7bce783404"),
    "event_id": ObjectId("60f67e52c0555f7bce783402"),
    "queue_name": "hello2"
} ]);
session.commitTransaction(); session.endSession();

// ----------------------------
// Collection structure for work_queues_config
// ----------------------------
db.getCollection("work_queues_config").drop();
db.createCollection("work_queues_config");

// ----------------------------
// Documents of work_queues_config
// ----------------------------
session = db.getMongo().startSession();
session.startTransaction();
db = session.getDatabase("msg_center");
db.getCollection("work_queues_config").insert([ {
    _id: ObjectId("6133354d7caeb90e436b9c93"),
    "event_id": ObjectId("613324477caeb90e436b9c92"),
    "queue_name": "workQueue"
} ]);
db.getCollection("work_queues_config").insert([ {
    _id: ObjectId("6133354d7caeb90e436b9c94"),
    "event_id": ObjectId("613324477caeb90e436b9c92"),
    "queue_name": "workQueue"
} ]);
session.commitTransaction(); session.endSession();
