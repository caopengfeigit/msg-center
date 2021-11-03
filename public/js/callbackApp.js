var app = new Vue(
    {
        router,
        el: '#app',
        data: {
            form: {
                ProjectName: '',
                EventKey: '',
                StartDate: '',
                EndDate: '',
            },
            date: [],
            tableData: [],
            multipleSelection: [],
            pageSize: 20,
            total: 0,
            currentPage: 1,
            dialog: false,
            title: '回调记录',
            formLabelWidth: '35%',
            logData: {
                ProjectName: '',
                EventKey: '',
                QueueName: '',
                RequestHost: '',
                RequestPath: '',
                RequestType: '',
                RequestIsJson: false,
                RequestRes: false,
                RequestStatus: '',
                RequestData: '',
                RequestError: '',
                RequestResponse: '',
                CreatedAtStr: '',
            }
        },
        created(){
            this.searchCallbackLogs()
        },
        methods: {
            //列表数据
            searchCallbackLogs(clear = 0) {
                if (clear) {
                    this.form.page = 1
                    this.currentPage = 1
                }
                this.form.ProjectName = this.form.ProjectName.replace(/[\u4E00-\u9FA5]/g,'')
                this.form.EventKey = this.form.EventKey.replace(/[\u4E00-\u9FA5]/g,'')
                if (this.date != null && this.date.length > 0) {
                    this.form.StartDate = this.date[0]
                    this.form.EndDate = this.date[1]
                } else {
                    this.form.StartDate = ''
                    this.form.EndDate = ''
                }
                axios.get('/api/get-callback-logs',{params:this.form}).then(response => {
                    this.tableData = response.data.list
                    this.total = response.data.pagination.Total
                    this.pageSize = response.data.pagination.PageSize
                });
            },

            //翻页
            currentChange(page) {
                this.form.page = page
                this.searchCallbackLogs()
            },

            //显示详细信息对话框
            showDialog(row) {
                this.logData.ProjectName = row.ProjectName
                this.logData.EventKey = row.EventKey
                this.logData.QueueName = row.QueueName
                this.logData.RequestHost = row.RequestHost
                this.logData.RequestPath = row.RequestPath
                this.logData.RequestType = row.RequestType
                this.logData.RequestIsJson = row.RequestIsJson
                this.logData.RequestRes = row.RequestRes
                this.logData.RequestStatus = row.RequestStatus
                this.logData.RequestData = row.RequestData
                this.logData.RequestError = row.RequestError
                this.logData.RequestResponse = row.RequestResponse
                this.logData.CreatedAtStr = row.CreatedAtStr
                this.dialog = true
            },

            //关闭添加业务线对话框
            closeDialog() {
                this.dialog = false
                this.logData.ProjectName = ''
                this.logData.EventKey = ''
                this.logData.QueueName = ''
                this.logData.RequestHost = ''
                this.logData.RequestPath = ''
                this.logData.RequestType = ''
                this.logData.RequestIsJson = false
                this.logData.RequestRes = false
                this.logData.RequestStatus = ''
                this.logData.RequestData = ''
                this.logData.RequestError = ''
                this.logData.RequestResponse = ''
                this.logData.CreatedAtStr = ''
            },
        }
    }
)