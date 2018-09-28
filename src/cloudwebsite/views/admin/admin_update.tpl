{{template "header.tpl" .}}

<script type="text/javascript">
    function checkData(){
        if(document.getElementById('account').value == ""){
            alert("账号不能为空!");
            return false;
        }
        return true;
    }
</script>

<div class="bg-light">
{{if .post}}
    {{if .error}}
    <div class="alert alert-danger">更新失败：{{str2html .error}}</div>
    {{else if .post}}
    <div class="alert alert-success">更新成功</div>
    <br/>
    {{end}}
{{end}}
  <form method="post" action="update" onsubmit="return checkData();">
    <div class="row">
        <div class="col">
        </div>
        <div class="col">
            <div class="row mt-2">
                <div class="col">
                    <div class="form-group">
                        <label for="account">账号：</label>
                        <input type="text" class="form-control" name="account" id="account" placeholder="" value="{{.admin.Account}}" readonly>
                    </div>
                </div>
                <div class="col">
                </div>
            </div>
            
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="adminadmin" name="adminadmin" {{if .admin.Adminadmin}}checked{{end}}>
                <label class="form-check-label" for="adminadmin">管理员管理</label>
            </div>
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="adminaccount" name="adminaccount" {{if .admin.Adminaccount}}checked{{end}}>
                <label class="form-check-label" for="adminaccount">用户管理</label>
            </div>
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="adminapp" name="adminapp" {{if .admin.Adminapp}}checked{{end}}>
                <label class="form-check-label" for="adminapp">应用管理</label>
            </div>
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="adminappdata" name="adminappdata" {{if .admin.Adminappdata}}checked{{end}}>
                <label class="form-check-label" for="adminappdata">应用数据管理</label>
            </div>
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="adminonline" name="adminonline" {{if .admin.Adminonline}}checked{{end}}>
                <label class="form-check-label" for="adminonline">在线用户管理</label>
            </div>
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="adminmessage" name="adminmessage" {{if .admin.Adminmessage}}checked{{end}}>
                <label class="form-check-label" for="adminmessage">用户消息管理</label>
            </div>
            <div class="row mt-2">
                <div class="col">
                    <div class="form-group">
                        <label for="expire">过期时间：</label>
                        <input type="text" class="form-control" id="expire" name="expire" value="{{dateformat .admin.Expire "01/02/2006"}}">
                    </div>   
                </div>
                <div class="col">
                </div>
            </div>
            <button type="submit" class="btn btn-outline-primary btn-lg btn-block mb-3" style="width:100px;">更新</button>
        </div>
        <div class="col">
        </div>
    </div>
  </form>
</div>

<script type="text/javascript">
  $( function() {
    $( "#expire" ).datepicker();
  } );
</script>

{{template "footer.tpl" .}}