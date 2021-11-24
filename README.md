### 下载地址
>
> [github](https://raw.githubusercontent.com/IoriLaw/IoriLaw.github.io/main/notice.zip)
>
> or
>
> [蓝奏云](https://wwp.lanzoui.com/iO4ijwup2od)


### 文件说明

| 文件名|备注|
| ------ | ----|
| dingDingNotification | 主程序|
| conf.json | 程序配置文件|
| notice.conf | nginx配置文件|


### 准备工作
>
>1. 打开防火墙端口 11223，程序写死不可变
>
>
>2. 程序和配置文件于同一目录，赋主程序执行权限
>
>
>3. nginx监听端口，通过include 配置文件方式， 或者根据配置文件自行配置
>
>
```
赋主程序执行权限
sudo chmod +x dingDingNotification
```
```
include 配置文件
http {
    ...
    include       /path_to_nginx_conf_dir/notice.conf
    ...
}

```


### 修改配置文件
```
{
  "Url": "https://oapi.dingtalk.com/robot/send?access_token=", //钉钉机器人地址
  "At": { // 需要在群里@对象属性配置
    "pc-创作者中心":"18888888888", //pc-创作者中心 项目变动时@18888888888手机用户
    "management-hub": "18888888888", //management-hub 项目变动时@18888888888手机用户
    ...
    "test-group": "18888888888",//test-group测试组，不限项目。变动时@18888888888手机用户
    ...
  },
  "keyword": "bbyy" // 简单秘钥，机器人创建时指定不可变
}
```
> 通常只需要增删对应组的手机号码。重启程序即可

| 项目|各端|备注|
| ------ | ----|---|
| pc-创作者中心<br> management-hub<br> creator-hub<br> pc-hub| 前端|
| backend-hub | 后端|
| ios-hub | iOS客户端|
| android-hub | Android客户端|
| test-group | 测试组| 测试能收到所有端。特殊组



### 启动程序
```
nohup ./dingDingNotification &
```

### 重启
```
ps -ef | grep dingDingNotification

然后kill掉， 再启动。
```


### 域名变更
> 1. 解析到目标服务器
>
>
> 2. 修改gitlab webhook里面域名部分， 其他不动

