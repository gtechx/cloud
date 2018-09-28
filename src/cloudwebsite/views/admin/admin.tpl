{{template "header.tpl" .}}

  <div class="bg-light">
    <div class="d-flex flex-wrap">
        <div class="p-1">
            <label for="account">账号：</label>
            <input type="text" class="rounded" style="width:100px" name="account" id="account" placeholder="">
        </div>
        <div class="p-1">
            <label for="adminadmin" class="checkbox-inline">
                <input type="checkbox" class="" name="adminadmin" id="adminadmin" placeholder="">adminadmin
            </label>
        </div>
        <div class="p-1">
            <label for="adminaccount" class="checkbox-inline">
                <input type="checkbox" class="" name="adminaccount" id="adminaccount" placeholder="">adminaccount
            </label>
        </div>
        <div class="p-1">
            <label for="adminapp" class="checkbox-inline">
                <input type="checkbox" class="" name="adminapp" id="adminapp" placeholder="">adminapp
            </label>
        </div>
        <div class="p-1">
            <label for="adminappdata" class="checkbox-inline">
                <input type="checkbox" class="" name="adminappdata" id="adminappdata" placeholder="">adminappdata
            </label>
        </div>
        <div class="p-1">
            <label for="adminonline" class="checkbox-inline">
                <input type="checkbox" class="" name="adminonline" id="adminonline" placeholder="">adminonline
            </label>
        </div>
        <div class="p-1">
            <label for="adminmessage" class="checkbox-inline">
                <input type="checkbox" class="" name="adminmessage" id="adminmessage" placeholder="">adminmessage
            </label>
        </div>
        <div class="p-1">
            <label for="expire">expire日期：</label>
            <input type="text" class="rounded" style="width:100px" name="expire" id="expire" placeholder="">
        </div>
        <div class="p-1">
            <button id="btn_filter" onclick="$('#table').bootstrapTable('refresh');" type="button" class="btn btn-info btn-sm">
                过滤
            </button>
        </div>
    </div>
    <div id="toolbar" class="btn-group">
        <button id="btn_add" onclick="window.location.href='create';" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-plus"></span>新增
        </button>
        <button id="btn_delete" onclick="delAdmin();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-x"></span>删除
        </button>
    </div>
    <table id="table">
    </table>
  </div>

<script type="text/javascript">
  $( function() {
    $( "#expire" ).datepicker();
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
  function delAdmin(){
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
              account: $("#account").val(),
              adminadmin: $("#adminadmin").prop("checked") ? "on" : "",
              adminaccount: $("#adminaccount").prop("checked") ? "on" : "",
              adminapp: $("#adminapp").prop("checked") ? "on" : "",
              adminappdata: $("#adminappdata").prop("checked") ? "on" : "",
              adminonline: $("#adminonline").prop("checked") ? "on" : "",
              adminmessage: $("#adminmessage").prop("checked") ? "on" : "",
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
                  if(value == {{.account}})
                  return value;
                  else
                  return '<a class="" href="update?account='+value+'">'+value+'</a>';
              }
          }, {
              field: 'adminadmin',
              title: 'adminadmin',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'adminaccount',
              title: 'adminaccount',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'adminapp',
              title: 'adminapp',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'adminappdata',
              title: 'adminappdata',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'adminonline',
              title: 'adminonline',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'adminmessage',
              title: 'adminmessage',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'expire',
              title: 'expire',
              align: 'center',
              valign: 'middle'
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