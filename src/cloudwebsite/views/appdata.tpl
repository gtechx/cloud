{{template "header.tpl" .}}

  <div class="bg-light">
    <div class="d-flex flex-wrap">
        <div class="p-1">
            <label for="appname">应用名称：</label>
            <select class="rounded" onchange="onAppnameChange(this);" style="width:100px" name="appname" id="appname">
                {{range $index, $elem := .applist}}
                <option>{{$elem.Appname}}</option>
                {{end}}
            </select>
        </div>
        <div class="p-1">
            <label for="zonename">分区名：</label>
            <select class="rounded" style="width:100px" name="zonename" id="zonename">
                <option></option>
            </select>
        </div>
        <div class="p-1">
            <label for="id">ID：</label>
            <input type="text" class="rounded" style="width:100px" name="id" id="id" placeholder="">
        </div>
        <div class="p-1">
            <label for="account">账号：</label>
            <input type="text" class="rounded" style="width:100px" name="account" id="account" placeholder="">
        </div>
        <div class="p-1">
            <label for="nickname">昵称：</label>
            <input type="text" class="rounded" style="width:100px" name="nickname" id="nickname" placeholder="">
        </div>
        <div class="p-1">
            <label for="sex">性别：</label>
            <select class="rounded" style="width:50px" name="sex" id="sex">
                <option></option>
                <option>男</option>
                <option>女</option>
            </select>
        </div>
        <div class="p-1">
            <label for="desc">描述：</label>
            <input type="text" class="rounded" style="width:150px" name="desc" id="desc" placeholder="">
        </div>
        <div class="p-1">
            <label for="email">Email：</label>
            <input type="text" class="rounded" style="width:100px" name="email" id="email" placeholder="">
        </div>
        <div class="p-1">
            <label for="ip">IP：</label>
            <input type="text" class="rounded" style="width:120px" name="ip" id="ip" placeholder="">
        </div>
        <div class="p-1">
            <label for="country">国家：</label>
            <input type="text" class="rounded" style="width:80px" name="country" id="country" placeholder="">
        </div>
        <div class="p-1">
            <label for="birthdaybegindate">生日起始：</label>
            <input type="text" class="rounded" style="width:100px" name="birthdaybegindate" id="birthdaybegindate" placeholder="">
        </div>
        <div class="p-1">
            <label for="birthdayenddate">生日最终：</label>
            <input type="text" class="rounded" style="width:100px" name="birthdayenddate" id="birthdayenddate" placeholder="">
        </div>
        <div class="p-1">
            <label for="lastloginbegindate">登录起始：</label>
            <input type="text" class="rounded" style="width:100px" name="lastloginbegindate" id="lastloginbegindate" placeholder="">
        </div>
        <div class="p-1">
            <label for="lastloginenddate">登录最终：</label>
            <input type="text" class="rounded" style="width:100px" name="lastloginenddate" id="lastloginenddate" placeholder="">
        </div>
        <div class="p-1">
            <label for="begindate">起始日期：</label>
            <input type="text" class="rounded" style="width:100px" name="begindate" id="begindate" placeholder="">
        </div>
        <div class="p-1">
            <label for="enddate">最终日期：</label>
            <input type="text" class="rounded" style="width:100px" name="enddate" id="enddate" placeholder="">
        </div>
        <div class="p-1">
        <button id="btn_filter" onclick="$('#table').bootstrapTable('refresh');" type="button" class="btn btn-info btn-sm">
            过滤
        </button>
        </div>
    </div>

    {{if not .isreadonly}}
    <div id="toolbar" class="btn-group">
        <button id="btn_add" onclick="window.location.href='create';" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-plus"></span>新增
        </button>
        <button id="btn_delete" onclick="delAppData();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-x"></span>删除
        </button>
        <button id="btn_ban" onclick="banAppDatas();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-ban"></span>封禁
        </button>
        <button id="btn_ban" onclick="unbanAppDatas();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-circle-check"></span>解除封禁
        </button>
    </div>
    {{end}}
    <table id="table">
    </table>
  </div>

<script type="text/javascript">
  $( function() {
    $( "#begindate" ).datepicker();
    $( "#enddate" ).datepicker();
    $( "#birthdaybegindate" ).datepicker();
    $( "#birthdayenddate" ).datepicker();
    $( "#lastloginbegindate" ).datepicker();
    $( "#lastloginenddate" ).datepicker();
    onAppnameChange(document.getElementsByName( "appname" )[0]);
  } );

  function onAppnameChange(obj) {
      if(obj.selectedIndex == -1)
        return;
      var opt = obj.options[obj.selectedIndex];
      console.info("text:"+opt.text);
      console.info("value:"+opt.value);

        $.post("zonelist", { 'account': $('#account').val(), 'appname': $('#appname').val()},
        function(data) {
        $('#zonename').bootstrapTable('refresh');
        var jsondata = JSON.parse(data);
        console.info(jsondata["rows"]);
        var liststr = '';
        var count = jsondata["total"];
        var html = "<option></option>";
        for (i in jsondata["rows"])
        {
            var row = jsondata["rows"][i];
            html += '<option>' + row.zonename + '</option>';
        }
        $('#zonename').html(html);
        });
  }

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

  function delAppData(){
    if (confirm("确认要删除吗？这将删除账号相关的所有数据！")==false){ 
        return; 
    }
    var selects = $('#table').bootstrapTable('getSelections');
    if(selects.length == 0)
        return;
    var strdata = new Array()
    for(i in selects){
      strdata[i] = selects[i].id
    }
    console.info("del appdata "+strdata)
    $.post("del", { 'appdata[]': strdata },
    function(data) {
      $('#table').bootstrapTable('refresh');
    });
  }

  function banAppDatas(){
    var selects = $('#table').bootstrapTable('getSelections');
    if(selects.length == 0)
        return;
    var strdata = new Array()
    for(i in selects){
      strdata[i] = selects[i].id
    }
    console.info("ban appdata "+strdata)
    $.post("ban", { 'appdata[]': strdata },
    function(data) {
      $('#table').bootstrapTable('refresh');
    });
  }

  function unbanAppDatas(){
    var selects = $('#table').bootstrapTable('getSelections');
    if(selects.length == 0)
        return;
    var strdata = new Array()
    for(i in selects){
      strdata[i] = selects[i].id
    }
    console.info("ban appdata "+strdata)
    $.post("unban", { 'appdata[]': strdata },
    function(data) {
      $('#table').bootstrapTable('refresh');
    });
  }

  function banAppData(index){
    var row = $('#table').bootstrapTable('getData')[index];

    var strdata = new Array()
    strdata[0] = row.id
    $.post("ban", { 'appdata[]': strdata },
    function(data) {
    $('#table').bootstrapTable('refresh');
    });
  }
  function unbanAppData(index){
    var row = $('#table').bootstrapTable('getData')[index];

    var strdata = new Array()
    strdata[0] = row.id
    $.post("unban", { 'appdata[]': strdata },
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
              id: $("#id").val(),
              account: $("#account").val(),
              appname: $("#appname").val(),
              zonename: $("#zonename").val(),
              nickname: $("#nickname").val(),
              sex: $("#sex").val(),
              desc: $("#desc").val(),
              email: $("#email").val(),
              ip: $("#ip").val(),
              country: $("#country").val(),
              birthdaybegindate: $("#birthdaybegindate").val(),
              birthdayenddate: $("#birthdayenddate").val(),
              lastloginbegindate: $("#lastloginbegindate").val(),
              lastloginenddate: $("#lastloginenddate").val(),
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
          {{if not .isreadonly}}
          {
              checkbox: true, // 显示一个勾选框
              align: 'center' // 居中显示
          },{{end}} {
              field: 'id',
              title: 'ID',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'account',
              title: '账号',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'appname',
              title: '应用名称',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'zonename',
              title: '分区名',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'nickname',
              title: '昵称',
              align: 'center',
              valign: 'middle',
              formatter: function (value, row, index) {
                  {{if .isreadonly}}
                  return value;
                  {{else}}
                  return '<a class="" href="update?id='+row.id+'">'+value+'</a>';
                  {{end}}
              }
          }, {
              field: 'desc',
              title: '描述',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'sex',
              title: '性别',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'birthday',
              title: '出生日期',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'country',
              title: '国家',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'regip',
              title: '注册ip',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'lastip',
              title: '上次登录ip',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'lastlogin',
              title: '上次登录日期',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'createdate',
              title: '注册日期',
              align: 'center',
              valign: 'middle'
          }{{if not .isreadonly}}, {
              field: 'isbaned',
              title: "操作",
              align: 'center',
              valign: 'middle',
              width: 160, // 定义列的宽度，单位为像素px
              formatter: function (value, row, index) {
                  if(value)
                    return '<button class="btn btn-primary btn-sm" onclick="unbanAppData('+index+');">解除封禁</button>';
                  else
                    return '<button class="btn btn-primary btn-sm" onclick="banAppData('+index+');">封禁</button>';
              }
          }{{end}}
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