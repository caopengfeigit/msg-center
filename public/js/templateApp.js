var app = new Vue(
    {
        router,
        el: '#app',
        data: {
            form: {
                ProjectName: '',
                EventName: '',
            },
            tableData: [],
            multipleSelection: [],
            pageSize: 20,
            total: 0,
            currentPage: 1,
            formLabelWidth: '35%',
            showTitlePrex: '事件：',
            showTitle: '',
            singleDialog: false,
            workQueuesDialog: false,
            publishSubscribeDialog: false,
            businessLine: '',
            eventType: '',
            exchangeType: '',
            singleData: {},
            workQueuesConfigList: [],
            publishSubscribeConfigList: [],
        },
        created(){
            this.searchConfigs()
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
            toggleSelection(rows) {
                if (rows) {
                    rows.forEach(row => {
                        this.$refs.multipleTable.toggleRowSelection(row);
                    });
                } else {
                    this.$refs.multipleTable.clearSelection();
                }
            },
            handleSelectionChange(val) {
                this.multipleSelection = val;
            },

            //列表数据
            searchConfigs(clear = 0) {
                if (clear) {
                    this.form.page = 1
                    this.currentPage = 1
                }
                this.form.ProjectName = this.form.ProjectName.replace(/[\u4E00-\u9FA5]/g,'')
                this.form.EventName = this.form.EventName.replace(/[\u4E00-\u9FA5]/g,'')
                axios.get('/api/get-list',{params:this.form}).then(response => {
                    this.tableData = response.data.list
                    this.total = response.data.pagination.Total
                    this.pageSize = response.data.pagination.PageSize
                });
            },

            //翻页
            currentChange(page) {
                this.form.page = page
                this.searchConfigs()
            },

            //删除配置请求
            delConfigRequest(id) {
                //删除配置
                axios.delete('/api/del-event-config',{params: {"id" : id}}).then(() => {
                    this.$message({
                        type: 'success',
                        message: '删除成功!'
                    });
                    this.searchConfigs()
                }).catch(() => {
                    this.$message({
                        type: 'error',
                        message: '删除失败!'
                    });
                });
            },

            //删除配置
            delConfig(eventId) {
                this.$confirm('确认是否要删除此配置?如果有未被消费的消息会被抛弃哟～', '提示', {
                    confirmButtonText: '确认',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    if (eventId) {
                        this.delConfigRequest(eventId)
                    }
                }).catch(() => {
                    this.$message({
                        type: 'info',
                        message: '已取消删除'
                    });
                });
            },

            //批量删除配置
            multiDelConfig() {
                if (this.multipleSelection.length == 0) {
                    this.$message({
                        type: 'info',
                        message: '请选择要删除的配置'
                    });
                } else {
                    ids = ''
                    for (var i = 0; i < this.multipleSelection.length; i++) {
                        ids += this.multipleSelection[i].EventId + ','
                    }
                    if (ids != '') {
                        this.$confirm('确认是否要删除所选配置?如果有未被消费的消息会被抛弃哟～', '提示', {
                            confirmButtonText: '确认',
                            cancelButtonText: '取消',
                            type: 'warning'
                        }).then(() => {
                            this.delConfigRequest(ids)
                        }).catch(() => {
                            this.$message({
                                type: 'info',
                                message: '已取消删除'
                            });
                        });
                    } else {
                        this.$message({
                            type: 'info',
                            message: '请选择要删除的模板'
                        });
                    }
                }
            },

            //查看详细配置信息
            showConfig(row) {
                this.showTitle = this.showTitlePrex + row.EventName
                this.businessLine = row.ProjectName
                this.eventType = row.EventType
                switch (row.EventType) {
                    case 'Single':
                        this.singleDialog = true
                        break
                    case 'WorkQueues':
                        this.workQueuesDialog = true
                        break
                    case 'PublishSubscribe':
                        this.exchangeType = row.ExchangeType
                        this.publishSubscribeDialog = true
                        break
                    default:
                        break
                }
                //获取该事件对应的配置数据（先展示对话框再加载数据）
                this.getConfigDetail(row.EventId, row.EventType)
            },

            //根据event id获取对应的配置数据
            getConfigDetail(eventId, eventType) {
                axios.get('/api/get-config-detail',{params: {'id': eventId}}).then(response => {
                    this.assignByEventType(eventType, response.data)
                }).catch(response => {
                    this.$message({
                        type: 'error',
                        message: response.message
                    });
                });
            },

            //根据不同的事件类型，进行数据赋值
            assignByEventType(eventType, data) {
                if (data != null) {
                    switch (eventType) {
                        case 'Single':
                            this.singleData = data
                            break
                        case "WorkQueues":
                            this.workQueuesConfigList = data
                            break
                        case "PublishSubscribe":
                            this.publishSubscribeConfigList = data
                            break
                    }
                }
            },

            //关闭详细配置信息
            closeConfig(eventType) {
                this.showTitle = ''
                this.businessLine = ''
                this.eventType = ''
                switch (eventType) {
                    case 'Single':
                        this.singleData = {}
                        this.singleDialog = false
                        break
                    case 'WorkQueues':
                        this.workQueuesConfigList = []
                        this.workQueuesDialog = false
                        break
                    case 'PublishSubscribe':
                        this.publishSubscribeConfigList = []
                        this.publishSubscribeDialog = false;
                        break
                    default:
                        break
                }
            },

            //添加配置跳转
            addConfig() {
                window.location.href = '/config/add-config'
            },
        }
    }
)