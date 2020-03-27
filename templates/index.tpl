<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>ScarletWAF Managent System</title>
</head>
<style>
    .middle {
        position: absolute;
        top: 0;
        right: 0;
        left: 0;
        bottom: 0;
        margin: auto auto;
        width: 400px;
        height: 500px;
        border-radius: 30px;
        backdrop-filter: blur(15px);
        display: flex;
        flex-direction: column;
    }

    body {
        background: url("https://picfiles.alphacoders.com/223/223042.jpg");
        object-fit: scale-down;
    }

    .register {
        border: 2px;
        border-radius: 8px;
        display: inline-block;
        margin: 4px auto;
    }
</style>
<body>
<div class="middle">
    <div style="margin-bottom: 8px;margin-top: 8px;">
        <p>{{ .WelcomeMsg }}</p>
    </div>
    <form action="/user/add" method="POST" style="margin-top:8px;">
        <div>
            <label for="username">用户名</label>
            <input type="text" name="name" id="username">
        </div>
        <div>
            <label for="Email">邮箱</label>
            <input type="email" name="email" id="Email">
        </div>
        <div>
            <label for="Password">密码</label>
            <input type="password" name="password" id="Password">
        </div>
        <div>
            <button class="register">注册</button>
        </div>
    </form>
</div>
</body>
</html>