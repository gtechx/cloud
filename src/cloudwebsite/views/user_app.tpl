{{template "header.tpl" .}}

<div class="row">
  <div class="col-2">
  </div>
  <div class="col-8 bg-light px-0">
      <ol class="breadcrumb">
        <li class="breadcrumb-item"><a href="index">主菜单</a></li>
        <li class="breadcrumb-item active" aria-current="page">我的应用</li>
      </ol>
  </div>
  <div class="col-2">
  </div>
</div>

<div class="row">
  <div class="col-2">
  </div>
  <div class="col-8 bg-light">
    <div id="toolbar" class="btn-group">
        <button id="btn_add" onclick="window.location.href='appcreate';" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-plus"></span>新增
        </button>
        <button id="btn_delete" onclick="delApp();" type="button" class="btn btn-info btn-sm rightSize">
            <span class="oi oi-x"></span>删除
        </button>
    </div>
    <table id="table">
    </table>
  </div>
  <div class="col-2">
  </div>
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
            <button type="submit" class="btn btn-outline-primary col-sm-2 col-form-label ml-2">添加</button>
          </div>   
        </form>
      </div>
    </div>
  </div>
</div>

<script type="text/javascript">
  $('#zonepanel').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget); // Button that triggered the modal
    var appname = button.data('whatever'); // Extract info from data-* attributes
    // If necessary, you could initiate an AJAX request here (and then do the updating in a callback).
    // Update the modal's content. We'll use jQuery here, but you could use a data binding library or other methods instead.
    var modal = $(this);
    modal.find('.modal-title').text(appname + ' 分区管理');
    $('#zoneappname').val(appname);
    $.post("zonelist", { 'appname': appname },
    function(data) {
      console.info("Data Loaded: " + data);
      var jsondata = JSON.parse(data);
      var liststr = '';
      for (i in jsondata)
      {
        var zone = jsondata[i]["name"];//.replace(/"/g, '');
        liststr += '<label class="checkbox-inline border border-success ml-2 bg-danger">\n';
        liststr += '<input type="checkbox" id="\''+zone+'\'" name="zoneitem" value="\''+zone+'\'">' + zone + '\n';
        liststr += '</label>';
      }
      $('#zonelist').html(liststr)
    });
  });
  function addZone(){
    $.post("zoneadd", { 'appname': $('#zoneappname').val(), 'zonename': $('#zonename').val() },
    function(data) {
      var jsondata = JSON.parse(data);
      var liststr = '';
      for (i in jsondata)
      {
        var zone = jsondata[i]["name"];//.replace(/"/g, '');
        liststr += '<label class="checkbox-inline border border-success ml-2 bg-danger">\n';
        liststr += '<input type="checkbox" id="\''+zone+'\'" name="zoneitem" value="\''+zone+'\'">' + zone + '\n';
        liststr += '</label>';
      }
      $('#zonelist').html(liststr)
    });
    return false;
  }

  function delZone(){
    obj = document.getElementsByName("zoneitem");
    check_val = [];
    for(k in obj){
        if(obj[k].checked)
            check_val.push(obj[k].value);
    }

    $.post("zonedel", { 'appname': $('#zoneappname').val(), 'zonename[]': check_val },
    function(data) {
      var jsondata = JSON.parse(data);
      if(jsondata["error"] != "")
        alert(jsondata["error"]);
      $('#table').bootstrapTable('refresh');
    });
  }

  function delApp(){
    if (confirm("确认要删除吗？这将删除应用相关的所有数据！")==false){ 
        return; 
    }
    var selects = $('#table').bootstrapTable('getSelections');
    var strdata = new Array()
    for(i in selects){
      strdata[i] = selects[i].name
    }
    console.info(strdata)
    $.post("appdel", { 'appname[]': strdata },
    function(data) {
      $('#table').bootstrapTable('refresh');
    });
  }

    $("#table").bootstrapTable({ // 对应table标签的id
      url: "applist", // 获取表格数据的url
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
                  return '<a class="" href="appmodify?appname='+value+'">'+value+'</a>';
              }
          }, {
              field: 'desc',
              title: '描述',
              align: 'center',
              valign: 'middle'
          }, {
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
                  return '<button class="btn btn-primary btn-sm" data-toggle="modal" data-target="#zonepanel" data-whatever="'+row.name+'">分区管理</button>';
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