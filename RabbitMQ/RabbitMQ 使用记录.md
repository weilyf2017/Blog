目前的exchange的路由策略是：每个需要队列的服务独享一个队列（queue），消费者（consumer）采用ACK自动应答模式处理队列消息。

如果需要新增一个队列服务，需要做如下开发步骤：

1.创建队列，发送消息


```
<?php

$routingkey = 'key';
//设置你的连接
$conn_args = array('host' => 'localhost', 'port' => '5672', 'login' => 'guest', 'password' => 'guest');
$conn = new AMQPConnection($conn_args);
if ($conn->connect()) {
    echo "Established a connection to the broker \n";
} else {
    echo "Cannot connect to the broker \n ";
}
//你的消息
$message = json_encode(array('Hello World3!', 'php3', 'c++3:'));
//创建channel
$channel = new AMQPChannel($conn);
//创建exchange
$ex = new AMQPExchange($channel);
$ex->setName('exchange2'); //创建名字
$ex->setType(AMQP_EX_TYPE_DIRECT);
$ex->setFlags(AMQP_DURABLE);
echo "exchange2 status:" . $ex->declareExchange();
echo "\n";
for ($i = 0; $i < 100; $i++) {
    if ($routingkey == 'key2') {
        $routingkey = 'key';
    } else {
        $routingkey = 'key2';
    }
    $ex->publish($message, $routingkey);
}
```

这样就产生了50条消息，但是没有消费者，所以没有被消费



2.创建消费者，消费信息


```
<?php

set_time_limit(0);
$e_name = 'exchange2'; //交换机名
$q_name = 'queue2'; //队列名
$k_route = 'key2'; //路由key 
//连接RabbitMQ
$conn_args = array('host' => '127.0.0.1', 'port' => '5672', 'login' => 'guest', 'password' => 'guest', 'vhost' => '/');
$conn = new AMQPConnection($conn_args);
$conn->connect();

$channel = new AMQPChannel($conn);   
//创建交换机
$ex = new AMQPExchange($channel);
$ex->setName($e_name);
$ex->setType(AMQP_EX_TYPE_DIRECT); //direct类型
$ex->setFlags(AMQP_DURABLE); //持久化
echo "Exchange Status:" . $ex->declareExchange() . "\n";

//创建队列
$q = new AMQPQueue($channel);
$q->setName($q_name);
$q->setFlags(AMQP_DURABLE); //持久化  
//绑定交换机与队列，并指定路由键
echo 'Queue Bind: ' . $q->bind($e_name, $k_route) . "\n"; //阻塞模式接收消息
echo "Message:\n";
$q->consume('processMessage', AMQP_AUTOACK); //自动ACK应答  
$conn->disconnect();

/** * 消费回调函数 * 处理消息 */
function processMessage($envelope, $queue) {
    var_dump($envelope->getRoutingKey());
    $msg = $envelope->getBody();
    echo $msg . "\n"; //处理消息
}
```

运行之后 可在 http://127.0.0.1:15672/  看到这个队列的详情情况

