<script type="text/javascript">
    function checkPassword(){
        if(document.getElementById('password1').value != document.getElementById('password2').value){
            alert("两次输入的密码不一致!");
            return false;
        }
        if(document.getElementById('password').value == "")
        {
            alert("密码不能为空!");
            return false;
        }
        return true;
    }
</script>

<form method="post" action="/main/register" onsubmit="return checkPassword();">
  <div class="form-group">
    <label for="account">账号：</label>
    <input type="text" class="form-control" name="account" id="account" placeholder="Account">
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
  <button type="submit" class="btn btn-outline-primary btn-lg btn-block">注册</button>
</form>
