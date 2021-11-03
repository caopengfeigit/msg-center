var app = new Vue(
    {
        router,
        el: '#app',
        data: {
            form: {
                projectId: '',
                projectName: '',
                businessLine: '',
                eventName: '',
                eventType: 'Single',
                workQueuesNum: 2,
                exchangeType: 'direct',
                exchangeName: '',
                routingKey: '',
                queueName: '',
                callbackHost: '',
                callbackPath: '',
                callbackRequestType: 'GET',
                callbackRequestIsJson: false,
                xMatch: 'any',
                fanoutConfigs: [
                    {
                        queueName: '',
                        callbackHost: '',
                        callbackPath: '',
                        callbackRequestType: 'GET',
                        callbackRequestIsJson: false,
                    },
                ],
                headers: [
                    {
                        key: '',
                        value: '',
                    },
                ],
                topic: [
                    {
                        routingKey: '',
                        queueName: '',
                        callbackHost: '',
                        callbackPath: '',
                        callbackRequestType: 'GET',
                        callbackRequestIsJson: false,
                    },
                ],
            },
            eventType: ['Single', 'WorkQueues', 'PublishSubscribe'],
            imgSrc: '/public/img/single.png',
            callbackRequestType: ['GET', 'POST'],
            exchangeType: {
                'direct': 'direct',
                'fanout': 'fanout',
                'headers': 'headers',
                'topic': 'topic',
                'x-delayed-message': 'x-delayed-message(延时队列)',
            },
            xMatch: ['any', 'all'],
            showSingle: true,
            showWorkqueues: false,
            showPublishSubscribe: false,
            showDirect: true,
            showFanout: false,
            showHeaders: false,
            showTopic: false,
            showXDelayedMessage: false,
        },
        created(){
            axios.interceptors.response.use(
                (response) => {
                    return response
                },
                (error) => {
                    return Promise.reject(error.response.data)
                }
            )
        },
        methods: {
            //搜索业务线
            searchProject() {
                this.form.projectName = this.form.projectName.replace(/[\u4E00-\u9FA5]/g,'')
                if (this.form.projectName == '') {
                    this.form.projectId = ''
                    this.$message({
                        type: 'error',
                        message: '请输入业务线标识!'
                    });
                    return
                }
                //请求接口获取业务线
                axios.get('/api/search-project',{params:{'projectName': this.form.projectName}}).then(response => {
                    if (!response.data.hasProject) {
                        this.form.projectId = ''
                        this.$message({
                            type: 'error',
                            message: '不存在该业务线!'
                        });
                    } else {
                        this.form.projectId = response.data.project.Id
                        this.form.businessLine = response.data.project.BusinessLine
                    }
                });
            },

            //搜索事件是否在当前业务线下存在
            searchEvent() {
                this.form.projectName = this.form.projectName.replace(/[\u4E00-\u9FA5]/g,'')
                this.form.eventName = this.form.eventName.replace(/[\u4E00-\u9FA5]/g,'')
                if (this.form.projectId == '') {
                    this.$message({
                        type: 'error',
                        message: '请选择业务线!'
                    });
                    return
                }
                if (this.form.eventName == '') {
                    this.$message({
                        type: 'error',
                        message: '请填写事件!'
                    });
                    return
                }

                //请求接口检查业务线是否被占用
                axios.get('/api/search-event',{params:{'projectId': this.form.projectId, 'eventName': this.form.eventName}}).then(response => {
                    if (response.data.hasEvent) {
                        this.form.eventName = ''
                        this.$message({
                            type: 'error',
                            message: '当前业务线下已存在该事件'
                        });
                    }
                }).catch(response => {
                    this.$message({
                        type: 'error',
                        message: response.message
                    });
                });
            },

            //事件类型改变
            changeType() {
                switch(this.form.eventType) {
                    case "Single":
                        this.imgSrc = '/public/img/single.png'
                        this.showWorkqueues = false
                        this.showPublishSubscribe = false
                        this.showSingle = true
                        break
                    case "WorkQueues":
                        this.imgSrc = '/public/img/workQueues.png'
                        this.showSingle = false
                        this.showPublishSubscribe = false
                        this.showWorkqueues = true
                        break
                    case "PublishSubscribe":
                        this.showSingle = false
                        this.showWorkqueues = false
                        this.showPublishSubscribe = true
                        this.changeExchangeType()
                        break
                }
            },

            //交换器类型改变
            changeExchangeType() {
                switch(this.form.exchangeType) {
                    case "direct":
                        this.imgSrc = '/public/img/direct.png'
                        this.showFanout = false
                        this.showHeaders = false
                        this.showTopic = false
                        this.showXDelayedMessage = false
                        this.showDirect = true
                        break
                    case "fanout":
                        this.imgSrc = '/public/img/fanout.png'
                        this.showDirect = false
                        this.showHeaders = false
                        this.showTopic = false
                        this.showXDelayedMessage = false
                        this.showFanout = true
                        break
                    case "headers":
                        this.imgSrc = '/public/img/headers1.png'
                        this.showDirect = false
                        this.showFanout = false
                        this.showTopic = false
                        this.showXDelayedMessage = false
                        this.showHeaders = true
                        break
                    case "topic":
                        this.imgSrc = '/public/img/topic.png'
                        this.showDirect = false
                        this.showFanout = false
                        this.showHeaders = false
                        this.showXDelayedMessage = false
                        this.showTopic = true
                        break
                    case "x-delayed-message":
                        this.imgSrc = '/public/img/delay.png'
                        this.showDirect = false
                        this.showFanout = false
                        this.showHeaders = false
                        this.showTopic = false
                        this.showXDelayedMessage = true
                        break
                }
            },

            //删除fanout配置
            delFanoutConfig(index) {
                if (this.form.fanoutConfigs.length - 1 < 1) {
                    this.$message({
                        type: 'error',
                        message: '至少设置一组配置'
                    });
                } else {
                    this.form.fanoutConfigs.splice(index, 1)
                }
            },

            //增加fanout配置
            addFanoutConfig() {
                this.form.fanoutConfigs.push({
                    queueName: '',
                    callbackHost: '',
                    callbackPath: '',
                    callbackRequestType: 'GET',
                    callbackRequestIsJson: false,
                })
            },

            //检测fanout配置中的队列名
            checkFanoutQueueName(index) {
                this.form.fanoutConfigs[index].queueName = this.form.fanoutConfigs[index].queueName.replace(/[\u4E00-\u9FA5]/g,'')
                if (this.form.fanoutConfigs[index].queueName != '') {
                    for (let i = 0; i < this.form.fanoutConfigs.length; i++) {
                        if (index != i && this.form.fanoutConfigs[i].queueName == this.form.fanoutConfigs[index].queueName) {
                            this.form.fanoutConfigs[index].queueName = ''
                            this.$message({
                                type: 'error',
                                message: '不能与其它队列名相同'
                            });
                            break;
                        }
                    }
                }
            },

            //删除header
            delHeader(index) {
                if (this.form.headers.length - 1 < 1) {
                    this.$message({
                        type: 'error',
                        message: '至少设置一组配置'
                    });
                } else {
                    this.form.headers.splice(index, 1)
                }
            },

            //增加header
            addHeader() {
                this.form.headers.push({
                    key: '',
                    value: '',
                })
            },

            //检查headerkey是否重复
            checkHeadersKey(index) {
                this.form.headers[index].key = this.form.headers[index].key.replace(/[\u4E00-\u9FA5]/g,'')
                if (this.form.headers[index].key != '') {
                    for (let i = 0; i < this.form.headers.length; i++) {
                        if (index != i && this.form.headers[i].key == this.form.headers[index].key) {
                            this.form.headers[index].key = ''
                            this.$message({
                                type: 'error',
                                message: '不能与其它key相同'
                            });
                            break;
                        }
                    }
                }
            },

            //删除topic
            delTopicConfig(index) {
                if (this.form.topic.length - 1 < 1) {
                    this.$message({
                        type: 'error',
                        message: '至少设置一组配置'
                    });
                } else {
                    this.form.topic.splice(index, 1)
                }
            },

            //检查topic routing key
            checkTopicRoutingKey(index) {
                this.form.topic[index].routingKey = this.form.topic[index].routingKey.replace(/[\u4E00-\u9FA5]/g,'')
                if (this.form.topic[index].routingKey != '') {
                    for (let i = 0; i < this.form.topic.length; i++) {
                        if (index != i && this.form.topic[i].routingKey == this.form.topic[index].routingKey) {
                            this.form.topic[index].routingKey = ''
                            this.$message({
                                type: 'error',
                                message: '不能与其它路径相同'
                            });
                            break;
                        }
                    }
                }
            },

            //增加topic配置
            addTopicConfig() {
                this.form.topic.push({
                    routingKey: '',
                    queueName: '',
                    callbackHost: '',
                    callbackPath: '',
                    callbackRequestType: 'GET',
                    callbackRequestIsJson: false,
                })
            },

            //提交
            onSubmit() {
                this.form.projectId = this.form.projectId.replace(/[\u4E00-\u9FA5]/g,'')
                this.form.eventName = this.form.eventName.replace(/[\u4E00-\u9FA5]/g,'')
                if (
                    this.form.projectId == ''
                    ||
                    this.form.eventName == ''
                    ||
                    $.inArray(this.form.eventType, this.eventType) < 0
                ) {
                    this.$message({
                        type: 'error',
                        message: '请填写必填项'
                    });
                    return
                }

                let valid = true
                //根据eventType走不同的逻辑
                switch (this.form.eventType) {
                    case "Single":
                        valid = this.checkSingleForm()
                        break
                    case "WorkQueues":
                        valid = this.checkWorkQueuesForm()
                        break
                    case "PublishSubscribe":
                        valid = this.checkPublishSubscribeForm()
                        break
                }
                if (!valid) {
                    this.$message({
                        type: 'error',
                        message: '请填写并检查参数的合法性'
                    });
                    return
                }

                //请求接口
                axios.post('/api/add-config', this.form).then(() => {
                    this.$message({
                        type: 'success',
                        message: '添加成功'
                    });
                    window.location.href = '/config/index'
                }).catch(response => {
                    this.$message({
                        type: 'error',
                        message: response.message
                    });
                });
            },

            //净化参数
            clearParams() {
                this.form.exchangeName = this.form.exchangeName.replace(/[\u4E00-\u9FA5]/g,'')
                this.form.routingKey = this.form.routingKey.replace(/[\u4E00-\u9FA5]/g,'')
                this.form.queueName = this.form.queueName.replace(/[\u4E00-\u9FA5]/g,'')
                this.form.callbackHost = this.form.callbackHost.replace(/[\u4E00-\u9FA5]/g,'')
                this.form.callbackPath = this.form.callbackPath.replace(/[\u4E00-\u9FA5]/g,'')
            },

            //single数据校验
            checkSingleForm() {
                let valid = true
                this.clearParams()
                if (
                    this.form.queueName == ''
                    ||
                    this.form.callbackHost == ''
                    ||
                    this.form.callbackPath == ''
                    ||
                    $.inArray(this.form.callbackRequestType, this.callbackRequestType) < 0
                ) {
                    valid = false
                }
                return valid
            },

            //workQueues数据校验
            checkWorkQueuesForm() {
                let valid = true
                this.clearParams()
                if (
                    this.form.workQueuesNum < 2
                    ||
                    this.form.workQueuesNum > 5
                    ||
                    this.form.queueName == ''
                    ||
                    this.form.callbackHost == ''
                    ||
                    this.form.callbackPath == ''
                    ||
                    $.inArray(this.form.callbackRequestType, this.callbackRequestType) < 0
                ) {
                    valid = false
                }
                return valid
            },

            //publish subscribe数据校验
            checkPublishSubscribeForm() {
                let valid = true
                if (!this.exchangeType.hasOwnProperty(this.form.exchangeType)) {
                    valid = false
                } else {
                    this.clearParams()
                    if (this.form.exchangeName == '') {
                        valid = false
                    } else {
                        switch (this.form.exchangeType) {
                            case 'direct':
                                if (
                                    this.form.routingKey == ''
                                    ||
                                    this.form.queueName == ''
                                    ||
                                    this.form.callbackHost == ''
                                    ||
                                    this.form.callbackPath == ''
                                    ||
                                    $.inArray(this.form.callbackRequestType, this.callbackRequestType) < 0
                                ) {
                                    valid = false
                                }
                                break
                            case 'fanout':
                                let queueNames = []
                                for (let key = 0; key < this.form.fanoutConfigs.length; key++) {
                                    this.form.fanoutConfigs[key].queueName = this.form.fanoutConfigs[key].queueName.replace(/[\u4E00-\u9FA5]/g,'')
                                    this.form.fanoutConfigs[key].callbackHost = this.form.fanoutConfigs[key].callbackHost.replace(/[\u4E00-\u9FA5]/g,'')
                                    this.form.fanoutConfigs[key].callbackPath = this.form.fanoutConfigs[key].callbackPath.replace(/[\u4E00-\u9FA5]/g,'')
                                    if (
                                        this.form.fanoutConfigs[key].queueName == ''
                                        ||
                                        this.form.fanoutConfigs[key].callbackHost == ''
                                        ||
                                        this.form.fanoutConfigs[key].callbackPath == ''
                                        ||
                                        $.inArray(this.form.fanoutConfigs[key].callbackRequestType, this.callbackRequestType) < 0
                                    ) {
                                        valid = false
                                        break
                                    }
                                    if ($.inArray(this.form.fanoutConfigs[key].queueName, queueNames) >= 0) {
                                        valid = false
                                        break
                                    }
                                    queueNames.push(this.form.fanoutConfigs[key].queueName)
                                }
                                break
                            case 'headers':
                                if (
                                    this.form.queueName == ''
                                    ||
                                    this.form.callbackHost == ''
                                    ||
                                    this.form.callbackPath == ''
                                    ||
                                    $.inArray(this.form.xMatch, this.xMatch) < 0
                                    ||
                                    $.inArray(this.form.callbackRequestType, this.callbackRequestType) < 0
                                ) {
                                    valid = false
                                } else {
                                    let headerKeys = []
                                    for (let key = 0; key < this.form.headers.length; key++) {
                                        this.form.headers[key].key = this.form.headers[key].key.replace(/[\u4E00-\u9FA5]/g,'')
                                        this.form.headers[key].value = this.form.headers[key].value.replace(/[\u4E00-\u9FA5]/g,'')
                                        if (
                                            this.form.headers[key].key == ''
                                            ||
                                            this.form.headers[key].value == ''
                                        ) {
                                            valid = false
                                            break
                                        }
                                        if ($.inArray(this.form.headers[key].key, headerKeys) >= 0) {
                                            valid = false
                                            break
                                        }
                                        headerKeys.push(this.form.headers[key].key)
                                    }
                                }
                                break
                            case 'topic':
                                let routingKeys = []
                                for (let key = 0; key < this.form.topic.length; key++) {
                                    this.form.topic[key].routingKey = this.form.topic[key].routingKey.replace(/[\u4E00-\u9FA5]/g,'')
                                    this.form.topic[key].queueName = this.form.topic[key].queueName.replace(/[\u4E00-\u9FA5]/g,'')
                                    this.form.topic[key].callbackHost = this.form.topic[key].callbackHost.replace(/[\u4E00-\u9FA5]/g,'')
                                    this.form.topic[key].callbackPath = this.form.topic[key].callbackPath.replace(/[\u4E00-\u9FA5]/g,'')
                                    if (
                                        this.form.topic[key].routingKey == ''
                                        ||
                                        this.form.topic[key].queueName == ''
                                        ||
                                        this.form.topic[key].callbackHost == ''
                                        ||
                                        this.form.topic[key].callbackPath == ''
                                        ||
                                        $.inArray(this.form.topic[key].callbackRequestType, this.callbackRequestType) < 0
                                    ) {
                                        valid = false
                                        break
                                    }
                                    if ($.inArray(this.form.topic[key].routingKey, routingKeys) >= 0) {
                                        valid = false
                                        break
                                    }
                                    routingKeys.push(this.form.topic[key].routingKey)
                                }
                                break
                            case 'x-delayed-message':
                                if (
                                    this.form.routingKey == ''
                                    ||
                                    this.form.queueName == ''
                                    ||
                                    this.form.callbackHost == ''
                                    ||
                                    this.form.callbackPath == ''
                                    ||
                                    $.inArray(this.form.callbackRequestType, this.callbackRequestType) < 0
                                ) {
                                    valid = false
                                }
                                break
                        }
                    }
                    return valid
                }
            },
        }
    }
)