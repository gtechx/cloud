{{template "header.tpl" .}}

  <div class="bg-light">
    <div class="d-flex flex-wrap">
        <div class="p-1">
            <label for="appnamefilter">应用名字：</label>
            <input type="text" class="rounded" style="width:100px" name="appnamefilter" id="appnamefilter" placeholder="">
        </div>
        <div class="p-1">
            <label for="descfilter">描述：</label>
            <input type="text" class="rounded" style="width:100px" name="descfilter" id="descfilter" placeholder="">
        </div>
        <div class="p-1">
            <label for="sharefilter">共享应用名字：</label>
            <input type="text" class="rounded" style="width:100px" name="sharefilter" id="sharefilter" placeholder="">
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

    <div id="toolbar" class="btn-group">
        <button id="btn_add" onclick="window.location.href='create';" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-plus"></span>新增
        </button>
        <button id="btn_delete" onclick="delApp();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-x"></span>删除
        </button>
    </div>
    <table id="table">
    </table>
  </div>

<div id="zonepanel" class="modal fade" tabindex="-1" role="dialog">
  <div class="modal-dialog modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title">Modal title</h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <div id="zonetoolbar" class="btn-group border-bottom mb-4 col-12">
          <button id="zone_delete" onclick="delZone();" type="button" class="btn btn-info btn-sm rightSize mb-2">
              <span class="oi oi-x"></span>批量删除
          </button>
        </div>
        <div class="d-flex flex-wrap" id="zonelist">
        </div>
      </div>
      <div class="modal-footer">
        <input type="text" class="invisible" id="zoneappname" name="zoneappname">
        <form onsubmit="return addZone();">
          <div class="form-group row">
            <label class="col-sm-5 col-form-label" for="zonename">分区名字：</label>
            <input type="text" class="form-control col-sm-4" id="zonename" name="zonename">
            <input type="hidden" class="form-control" id="account" name="account">
            <button type="submit" class="btn btn-outline-primary col-sm-2 col-form-label ml-2">添加</button>
          </div>   
        </form>
      </div>
    </div>
  </div>
</div>

<script type="text/javascript">
  function refreshzones(appname) {
    $.post("../zone/list", { 'appname': appname },
    function(data) {
      console.info("Data Loaded: " + data);
      var jsondata = JSON.parse(data);
      var liststr = '';
      for (i in jsondata)
      {
        var zone = jsondata[i]["zonename"];//.replace(/"/g, '');
        liststr += '<label class="checkbox-inline border border-success ml-2 bg-danger">\n';
        liststr += '<input type="checkbox" id="'+zone+'" name="zoneitem" value="'+zone+'">' + zone + '\n';
        liststr += '</label>';
      }
      $('#zonelist').html(liststr)
    });
  }
  $('#zonepanel').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget); // Button that triggered the modal
    var index = button.data('whatever'); // Extract info from data-* attributes
    var row = $('#table').bootstrapTable('getData')[index];
    // If necessary, you could initiate an AJAX request here (and then do the updating in a callback).
    // Update the modal's content. We'll use jQuery here, but you could use a data binding library or other methods instead.
    var modal = $(this);
    modal.find('.modal-title').text(row.appname + ' 分区管理');
    $('#zoneappname').val(row.appname);
    $('#account').attr("value", row.owner);
    $.post("../zone/list", { 'appname': row.appname },
    function(data) {
      console.info("Data Loaded: " + data);
      var jsondata = JSON.parse(data);
      var liststr = '';
      for (i in jsondata)
      {
        var zone = jsondata[i]["zonename"];//.replace(/"/g, '');
        liststr += '<label class="checkbox-inline border border-success ml-2 bg-danger">\n';
        liststr += '<input type="checkbox" id="'+zone+'" name="zoneitem" value="'+zone+'">' + zone + '\n';
        liststr += '</label>';
      }
      $('#zonelist').html(liststr)
    });
  });
  function addZone(){
    $.post("../zone/create", { 'appname': $('#zoneappname').val(), 'zonename': $('#zonename').val(), 'account': $('#account').val() },
    function(data) {
      var jsondata = JSON.parse(data);
      var liststr = '';
      for (i in jsondata)
      {
        var zone = jsondata[i]["zonename"];//.replace(/"/g, '');
        liststr += '<label class="checkbox-inline border border-success ml-2 bg-danger">\n';
        liststr += '<input type="checkbox" id="'+zone+'" name="zoneitem" value="'+zone+'">' + zone + '\n';
        liststr += '</label>';
      }
      $('#zonelist').html(liststr)
    });
    return false;
  }

  function delZone(){
    obj = document.getElementsByName("zoneitem");
    check_val = [];
    console.info(obj)
    for(var i = 0; i < obj.length; i++){
      console.info(obj[i].checked)
        if(obj[i].checked)
            check_val.push(obj[i].value);
    }
    
    $.post("../zone/del", { 'appname': $('#zoneappname').val(), 'zonename[]': check_val },
    function(data) {
      console.info(data)
      var jsondata = JSON.parse(data);
      if(jsondata["error"] != "")
        alert(jsondata["error"]);
      else
        refreshzones($('#zoneappname').val())
    });
  }

  function delApp(){
    if (confirm("确认要删除吗？这将删除应用相关的所有数据！")==false){ 
        return; 
    }
    var selects = $('#table').bootstrapTable('getSelections');
    var strdata = new Array()
    for(var i = 0; i < selects.length; i++){
      strdata.push(selects[i].appname);
    }
    $.post("del", { 'appname[]': strdata },
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
      pageList: [10, 15], // 设置页面可以显示的数据条数
      pageSize: 10, // 页面数据条数
      pageNumber: 1, // 首页页码
      sidePagination: 'server', // 设置为服务器端分页
      queryParamsType: "",
      queryParams: function (params) { // 请求服务器数据时发送的参数，可以在这里添加额外的查询参数，返回false则终止请求

          return {
              appnamefilter: $("#appnamefilter").val(),
              descfilter: $("#descfilter").val(),
              sharefilter: $("#sharefilter").val(),
              createbegindate: $("#begindate").val(),
              createenddate: $("#enddate").val(),
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
              field: 'appname',
              title: '名称',
              align: 'center',
              valign: 'middle',
              formatter: function (value, row, index) {
                  return '<a class="" href="update?appname='+value+'">'+value+'</a>';
              }
          }, {
              field: 'desc',
              title: '描述',
              align: 'center',
              valign: 'middle'
          }, {{if .isadmin}}{
              field: 'owner',
              title: '拥有者',
              align: 'center',
              valign: 'middle'
          },{{end}} {
              field: 'createdate',
              title: '创建日期',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'share',
              title: '共享数据应用名',
              align: 'center',
              valign: 'middle'
          }, {
              field: 'share',
              title: "操作",
              align: 'center',
              valign: 'middle',
              width: 160, // 定义列的宽度，单位为像素px
              formatter: function (value, row, index) {
                if(value == "")
                  return '<button class="btn btn-primary btn-sm" data-toggle="modal" data-target="#zonepanel" data-whatever="'+index+'">分区管理</button>';
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