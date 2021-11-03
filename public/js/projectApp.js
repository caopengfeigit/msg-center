var app = new Vue(
    {
        router,
        el: '#app',
        data: {
            form: {
                Name: '',
                BusinessLine: '',
            },
            tableData: [],
            multipleSelection: [],
            pageSize: 20,
            total: 0,
            currentPage: 1,
            formLabelWidth: '35%',
            addTableName: '新增业务线',
            editTableName: '编辑业务线',
            dialog: false,
            editDialog: false,
            addForm: {
                name: '',
                business_line: '',
            },
            editForm: {
                id: '',
                name: '',
                business_line: '',
            },
        },
        created(){
            this.searchProjects()
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
            searchProjects(clear = 0) {
                if (clear) {
                    this.form.page = 1
                    this.currentPage = 1
                }
                this.form.Name = this.form.Name.replace(/[\u4E00-\u9FA5]/g,'')
                axios.get('/api/get-project-list',{params:this.form}).then(response => {
                    this.tableData = response.data.list
                    this.total = response.data.pagination.Total
                    this.pageSize = response.data.pagination.PageSize
                });
            },

            //翻页
            currentChange(page) {
                this.form.page = page
                this.searchProjects()
            },

            //删除业务线请求
            delProjectRequest(id) {
                //删除业务线
                axios.delete('/api/del-project',{params: {"id" : id}}).then(() => {
                    this.$message({
                        type: 'success',
                        message: '删除成功!'
                    });
                    this.searchProjects()
                }).catch(() => {
                    this.$message({
                        type: 'error',
                        message: '删除失败!'
                    });
                });
            },

            //删除业务线
            delProject(Id) {
                this.$confirm('确认是否要删除此业务线？其下对应的所有配置以及未被消费的消息也将一并删除哟～', '提示', {
                    confirmButtonText: '确认',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    if (Id) {
                        this.delProjectRequest(Id)
                    }
                }).catch(() => {
                    this.$message({
                        type: 'info',
                        message: '已取消删除'
                    });
                });
            },

            //批量删除业务线
            multiDelProject() {
                if (this.multipleSelection.length == 0) {
                    this.$message({
                        type: 'info',
                        message: '请选择要删除的业务线'
                    });
                } else {
                    ids = ''
                    for (var i = 0; i < this.multipleSelection.length; i++) {
                        ids += this.multipleSelection[i].Id + ','
                    }
                    if (ids != '') {
                        this.$confirm('确认是否要删除所选业务线？其下对应的所有配置也将一并删除哟～', '提示', {
                            confirmButtonText: '确认',
                            cancelButtonText: '取消',
                            type: 'warning'
                        }).then(() => {
                            this.delProjectRequest(ids)
                        }).catch(() => {
                            this.$message({
                                type: 'info',
                                message: '已取消删除'
                            });
                        });
                    } else {
                        this.$message({
                            type: 'info',
                            message: '请选择要删除的业务线'
                        });
                    }
                }
            },

            //显示添加业务线对话框
            showProjectDialog() {
                this.dialog = true
            },

            //关闭添加业务线对话框
            closeDialog() {
                this.dialog = false
                this.addForm.name = ''
                this.addForm.business_line = ''
            },

            //添加业务线
            addProject() {
                this.addForm.name = this.addForm.name.replace(/[\u4E00-\u9FA5]/g,'')
                if (this.addForm.name == '' || this.addForm.business_line == '') {
                    this.$message({
                        type: 'info',
                        message: '请填写必填项'
                    });
                } else {
                    axios.post('/api/add-project', Qs.stringify(this.addForm)).then(() => {
                        this.$message({
                            type: 'success',
                            message: '添加成功'
                        });
                        this.dialog = false
                        this.searchProjects()
                        this.addForm.name = ''
                        this.addForm.business_line = ''
                    }).catch(response => {
                        this.$message({
                            type: 'error',
                            message: response.message
                        });
                    });
                }
            },

            //显示修改业务线对话框
            showEditProjectDialog(id, name, businessLine) {
                this.editForm.id = id
                this.editForm.name = name
                this.editForm.business_line = businessLine
                this.editDialog = true
            },

            //修改业务线
            editProject() {
                this.editForm.name = this.editForm.name.replace(/[\u4E00-\u9FA5]/g,'')
                if (this.editForm.id == '' || this.editForm.name == '' || this.editForm.business_line == '') {
                    this.$message({
                        type: 'info',
                        message: '请填写必填项'
                    });
                } else {
                    this.$confirm('确认是否要修改业务线？其下对应的所有配置将重启，且可能会导致消息丢失！', '提示', {
                        confirmButtonText: '确认',
                        cancelButtonText: '取消',
                        type: 'warning'
                    }).then(() => {
                        const loading = this.$loading({
                            lock: true,
                            text: 'Waiting...',
                            background: 'rgba(0, 0, 0, 0.7)'
                        });
                        axios.post('/api/edit-project', Qs.stringify(this.editForm)).then(() => {
                            this.$message({
                                type: 'success',
                                message: '修改成功'
                            });
                            loading.close();
                            this.editDialog = false
                            this.searchProjects()
                            this.editForm.id = ''
                            this.editForm.name = ''
                            this.editForm.business_line = ''
                        }).catch(response => {
                            loading.close();
                            this.$message({
                                type: 'error',
                                message: response.message
                            });
                        });
                    }).catch(() => {
                        this.$message({
                            type: 'info',
                            message: '已取消修改'
                        });
                    });
                }
            },
        }
    }
)