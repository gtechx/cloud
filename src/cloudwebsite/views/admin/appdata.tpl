{{template "header.tpl" .}}

  <div class="bg-light">
    <form class="form-inline">
        <div class="form-group col-2">
            <label for="accountfilter">账号：</label>
            <input type="text" class="form-control col-8" name="accountfilter" id="accountfilter" placeholder="">
        </div>
        <div class="form-group col-2">
            <label for="emailfilter">Email：</label>
            <input type="text" class="form-control col-8" name="emailfilter" id="emailfilter" placeholder="">
        </div>
        <div class="form-group col-2">
            <label for="ipfilter">IP：</label>
            <input type="text" class="form-control col-8" name="ipfilter" id="ipfilter" placeholder="">
        </div>
        <div class="form-group">
            <label for="begindate">起始日期：</label>
            <input type="text" class="form-control col-8" name="begindate" id="begindate" placeholder="">
        </div>
        <div class="form-group">
            <label for="enddate">最终日期：</label>
            <input type="text" class="form-control col-8" name="enddate" id="enddate" placeholder="">
        </div>
        <button id="btn_filter" onclick="$('#table').bootstrapTable('refresh');" type="button" class="btn btn-info btn-sm rightSize">
            过滤
        </button>
    </form>
    <div id="toolbar" class="btn-group">
        <button id="btn_add" onclick="addAccount();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-plus"></span>新增
        </button>
        <button id="btn_delete" onclick="delAccount();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-x"></span>删除
        </button>
        <button id="btn_ban" onclick="banAccounts();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-ban"></span>封禁
        </button>
        <button id="btn_ban" onclick="unbanAccounts();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-circle-check"></span>解除封禁
        </button>
    </div>
    <table id="table">
    </table>
  </div>

<div id="accountpanel" class="modal fade" tabindex="-1" role="dialog">
  <div class="modal-dialog modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title">Modal title</h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form>
            <div class="form-group">
                <label for="caccount">账号：</label>
                <input type="text" class="form-control" name="caccount" id="caccount" placeholder="Account">
            </div>
            <div class="form-group">
                <label for="cemail">邮箱：</label>
                <input type="email" class="form-control" name="cemail" id="cemail" placeholder="Email">
            </div>
            <div class="form-group">
                <label for="password1">密码：</label>
                <input type="password" class="form-control" name="password1" id="password1" placeholder="Password" oninput="document.getElementById('password').value = this.value;" onpropertychange="document.getElementById('password').value = this.value;">
                <input type="hidden" name="password" id="password" />
            </div>
            <div class="form-group">
                <label for="password2">确认密码：</label>
                <input type="password" class="form-control" name="password2" id="password2" placeholder="Password">
            </div>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-primary" id="createaccount" onclick="createAccount();">创建</button>
        <button type="button" class="btn btn-primary" id="modifyaccount" onclick="updateAccount();">修改</button>
        <button type="button" class="btn btn-secondary" data-dismiss="modal">取消</button>
      </div>
    </div>
  </div>
</div>

<script type="text/javascript">
  $( function() {
    $( "#begindate" ).datepicker();
  } );

  $( function() {
    $( "#enddate" ).datepicker();
  } );

  function checkEmail(email) {
    var reg = new RegExp("^[a-z0-9]+([._\\-]*[a-z0-9])*@([a-z0-9]+[-a-z0-9]*[a-z0-9]+.){1,63}[a-z0-9]+$"); //正则表达式

    if(email === ""){ //输入不能为空
        alert("邮箱不能为空!");
　　　　return false;
　　}else if(!reg.test(email)){ //正则验证不通过，格式不对
　　　　alert("邮箱格式不对!");
　　　　return false;
　　}
    return true;
  }

  function modifyAccount(index) {
        $( "#password" ).val("")
        $( "#password1" ).val("")
        $( "#password2" ).val("")
      $( "#createaccount" ).hide();
      $( "#modifyaccount" ).show();
      var row = $('#table').bootstrapTable('getData')[index];
      var modal = $('#accountpanel');
      modal.find('.modal-title').text("修改-"+row.account);
      $( "#caccount" ).val(row.account)
      $( "#cemail" ).val(row.email)
      modal.modal('show');
  }

  function updateAccount(){
    if(document.getElementById('password1').value != document.getElementById('password2').value){
        alert("两次输入的密码不一致!");
        return;
    }
    $.post("update", { 'account': $('#caccount').val(), 'email': $('#cemail').val(), 'password': $('#password').val() },
    function(data) {
    $('#table').bootstrapTable('refresh');
    });
  }

  function addAccount() {
    $( "#createaccount" ).show();
    $( "#modifyaccount" ).hide();
    var modal = $('#accountpanel');
    modal.find('.modal-title').text("创建");
    $( "#caccount" ).val("")
    $( "#cemail" ).val("")
    $( "#password" ).val("")
    $( "#password1" ).val("")
    $( "#password2" ).val("")
    modal.modal('show');
  }

  function createAccount(){
    if(document.getElementById('password1').value != document.getElementById('password2').value){
        alert("两次输入的密码不一致!");
        return;
    }
    if(document.getElementById('password').value == "")
    {
        alert("密码不能为空!");
        return;
    }
    $.post("create", { 'account': $('#caccount').val(), 'email': $('#cemail').val(), 'password': $('#password').val() },
    function(data) {
    $('#table').bootstrapTable('refresh');
    });
  }
  function delAccount(){
    if (confirm("确认要删除吗？这将删除账号相关的所有数据！")==false){ 
        return; 
    }
    var selects = $('#table').bootstrapTable('getSelections');
    if(selects.length == 0)
        return;
    var strdata = new Array()
    for(i in selects){
      strdata[i] = selects[i].account
    }
    console.info("del account "+strdata)
    $.post("del", { 'account[]': strdata },
    function(data) {
      $('#table').bootstrapTable('refresh');
    });
  }

  function banAccounts(){
    var selects = $('#table').bootstrapTable('getSelections');
    if(selects.length == 0)
        return;
    var strdata = new Array()
    for(i in selects){
      strdata[i] = selects[i].account
    }
    console.info("ban account "+strdata)
    $.post("ban", { 'account[]': strdata },
    function(data) {
      $('#table').bootstrapTable('refresh');
    });
  }

  function unbanAccounts(){
    var selects = $('#table').bootstrapTable('getSelections');
    if(selects.length == 0)
        return;
    var strdata = new Array()
    for(i in selects){
      strdata[i] = selects[i].account
    }
    console.info("ban account "+strdata)
    $.post("unban", { 'account[]': strdata },
    function(data) {
      $('#table').bootstrapTable('refresh');
    });
  }

  function banAccount(index){
    var row = $('#table').bootstrapTable('getData')[index];

    var strdata = new Array()
    strdata[0] = row.account
    $.post("ban", { 'account[]': strdata },
    function(data) {
    $('#table').bootstrapTable('refresh');
    });
  }
  function unbanAccount(index){
    var row = $('#table').bootstrapTable('getData')[index];

    var strdata = new Array()
    strdata[0] = row.account
    $.post("unban", { 'account[]': strdata },
    function(data) {
    $('#table').bootstrapTable('refresh');
    });
  }

    $("#table").bootstrapTable({ // 对应table标签的id
      url: "list", // 获取表格数据的url
      cache: false, // 设置为 false 禁用 AJAX 数据缓存， 默认为true
      clickToSelect: true,
      pagination: true,
      height: 650,
      toolbar: "#toolbar",
      striped: true,  //表格显示条纹，默认为false
      pagination: true, // 在表格底部显示分页组件，默认false
      pageList: [20, 50, 100], // 设置页面可以显示的数据条数
      pageSize: 20, // 页面数据条数
      pageNumber: 1, // 首页页码
      sidePagination: 'server', // 设置为服务器端分页
      queryParamsType: "",
      queryParams: function (params) { // 请求服务器数据时发送的参数，可以在这里添加额外的查询参数，返回false则终止请求

          return {
              accountfilter: $("#accountfilter").val(),
              emailfilter: $("#emailfilter").val(),
              ipfilter: $("#ipfilter").val(),
              begindate: $("#begindate").val(),
              enddate: $("#enddate").val(),
              pageSize: params.pageSize, // 每页要显示的数据条数
              //offset: params.offset, // 每页显示数据的开始行号
              pageNumber: params.pageNumber
              //sort: params.sort, // 要排序的字段
              //sortOrder: params.order, // 排序规则
              //dataId: $("#dataId").val() // 额外添加的参数
          }
      },
      //sortName: 'id', // 要排序的字段
      //sortOrder: 'desc', // 排序规则
      columns: [
          {
              checkbox: true, // 显示一个勾选框
              align: 'center' // 居中显示
          }, {
              field: 'account',
              title: '账号',
              align: 'center',
              valign: 'middle',
              formatter: function (value, row, index) {
                  return '<a class="" href="#" onclick="modifyAccount('+index+');">'+value+'</a>';
              }
          }, {
              field: 'email',
              title: '邮箱',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'regip',
              title: '注册ip',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'createdate',
              title: '注册日期',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'isbaned',
              title: "操作",
              align: 'center',
              valign: 'middle',
              width: 160, // 定义列的宽度，单位为像素px
              formatter: function (value, row, index) {
                  if(value)
                    return '<button class="btn btn-primary btn-sm" onclick="unbanAccount('+index+');">解除封禁</button>';
                  else
                    return '<button class="btn btn-primary btn-sm" onclick="banAccount('+index+');">封禁</button>';
              }
          }
      ],
      onLoadSuccess: function(data){  //加载成功时执行
            console.info("加载成功");
            //console.info(data);
      },
      onLoadError: function(status, res){  //加载失败时执行
            console.info("加载数据失败");
      }

});
</script>

{{template "footer.tpl" .}}