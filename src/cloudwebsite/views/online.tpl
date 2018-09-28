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
            <label for="country">国家：</label>
            <input type="text" class="rounded" style="width:80px" name="country" id="country" placeholder="">
        </div>
        <div class="p-1">
            <label for="platform">平台：</label>
            <select class="rounded" style="width:50px" name="platform" id="platform">
                <option></option>
                <option>web</option>
                <option>web1</option>
                <option>web2</option>
            </select>
        </div>
        <div class="p-1">
            <label for="serveraddr">服务器：</label>
            <select class="rounded" style="width:50px" name="serveraddr" id="serveraddr">
                <option></option>
                <option>127.0.0.1</option>
            </select>
        </div>
        <div class="p-1">
            <label for="onlinebegindate">登录起始日期：</label>
            <input type="text" class="rounded" style="width:100px" name="onlinebegindate" id="onlinebegindate" placeholder="">
        </div>
        <div class="p-1">
            <label for="onlineenddate">登录最终日期：</label>
            <input type="text" class="rounded" style="width:100px" name="onlineenddate" id="onlineenddate" placeholder="">
        </div>
        <div class="p-1">
        <button id="btn_filter" onclick="$('#table').bootstrapTable('refresh');" type="button" class="btn btn-info btn-sm">
            过滤
        </button>
        </div>
    </div>

    <div id="toolbar" class="btn-group">
        <button id="btn_ban" onclick="banAppDatas();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-ban"></span>封禁
        </button>
        <button id="btn_ban" onclick="banAppDatas();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-ban"></span>禁言
        </button>
        <button id="btn_ban" onclick="unbanAppDatas();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-circle-check"></span>解除禁言
        </button>
    </div>
    <table id="table">
    </table>
  </div>

<script type="text/javascript">
  $( function() {
    $( "#onlinebegindate" ).datepicker();
    $( "#onlineenddate" ).datepicker();
    //$("#country").countrySelect();
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
              country: $("#country").val(),
              platform: $("#platform").val(),
              serveraddr: $("#serveraddr").val(),
              begindate: $("#onlinebegindate").val(),
              enddate: $("#onlineenddate").val(),
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
              field: 'dataid',
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
              valign: 'middle'
          }, {
              field: 'sex',
              title: '性别',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'country',
              title: '国家',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'platform',
              title: '平台',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'serveraddr',
              title: '服务器',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'createdate',
              title: '登录日期',
              align: 'center',
              valign: 'middle'
          }{{if not .isreadonly}}, {
              field: 'isbaned',
              title: "操作",
              align: 'center',
              valign: 'middle',
              width: 160, // 定义列的宽度，单位为像素px
              formatter: function (value, row, index) {
                  var html = '';
                  if(!value)
                    html += '<button class="btn btn-primary btn-sm" onclick="banAppData('+index+');">封禁</button>';
                if(row.isjinyan)
                    html += '<button class="btn btn-primary btn-sm" onclick="unbanAppData('+index+');">解除禁言</button>';
                  else
                    html += '<button class="btn btn-primary btn-sm" onclick="banAppData('+index+');">禁言</button>';

                return html;
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