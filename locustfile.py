from random import randint
from locust import HttpLocust, TaskSet, task
import time
import itertools
import json
import uuid

runCounter = itertools.count()

def createOrderBody():
  return """ {
    "number": "1",
    "reference": "ref-001",
    "notes": "Nota",
    "price": 15000
    }"""

def payOrderBody(amount):
  return """ {
    "external_id": "10",
    "amount": amount,
    "type": "PAYMENT",
    "card_brand": "VISA",
    "card_bin": "1402",
    "card_last": "3211"
    }"""       

def createOrderItemBody(quantity):
  sku = uuid.uuid1()
  return """ {
    "sku": "%s",
    "unit_price": 1000,
    "quantity": %d,
    }""" % (sku, quantity)

def getTimeInMillis():
  return int(round(time.time() * 1000))

def getRandomTime():
  return randint(100,200) / 1000.0

def createOrderItem(l, orderID, quantity, currentRun):
  start = getTimeInMillis()
  orderItemUrl = "/orders/%s/items" % (orderID)
  createOrderItemResponse = l.client.post(orderItemUrl, createOrderItemBody(quantity))
  end = getTimeInMillis() - start
  print "createOrderItem, responseCode[%d], time[%dms], run[%d]" % (createOrderItemResponse.status_code, end, currentRun)  
  time.sleep(getRandomTime())

def payOrder(l, orderID, amount, currentRun):
  start = getTimeInMillis()
  orderItemUrl = "/orders/%s/transactions" % (orderID)
  payOrderResponse = l.client.post(orderItemUrl, payOrderBody(amount))
  end = getTimeInMillis() - start
  print "payOrder, responseCode[%d], time[%dms], run[%d]" % (payOrderResponse.status_code, end, currentRun)

def getOrder(l, orderID, currentRun):
  start = getTimeInMillis()
  getOrderUrl = "/orders/%s" % (orderID)
  getOrderResponse = l.client.get(getOrderUrl)
  end = getTimeInMillis() - start
  print "getOrder, responseCode[%d], time[%dms], run[%d]" % (getOrderResponse.status_code, end, currentRun)    

class RedWeddingBehaviour(TaskSet):
  def init(self, *args, **kwargs):
    super(HttpLocust, self).__init__(*args, **kwargs)
    self.orderID = None
    self.currentRun = None

  def on_start(self):
    self.createOrder()
  
  @task(1)
  def createOrder(self):
    self.currentRun = runCounter.next()
    #Cria a Ordem
    startCreateOrder = getTimeInMillis()
    createOrderResponse = self.client.post("/orders", createOrderBody())
    endCreateOrder = getTimeInMillis() - startCreateOrder
    print "createOrder, responseCode[%d], time[%dms], run[%d]" % (createOrderResponse.status_code, endCreateOrder, self.currentRun)
    if (createOrderResponse.status_code == 200):
      self.orderID = createOrderResponse.json()['id']

  @task(1)
  def addOrderItems(self):
    if (self.orderID != None):
      #5 consultas simultaneas--
      #Adiciona 5 itens de um produto a ordem, um OrderItem com o valor 5 na quantidade
      createOrderItem(self, self.orderID, 5, self.currentRun)
      #adiciona novamente um OrderItem com 5 unidades do produto
      createOrderItem(self, self.orderID, 5, self.currentRun)
      #adiciona um produto com somente uma unidade
      createOrderItem(self, self.orderID, 1, self.currentRun) #110.00 -> 50 + 50 + 10
      #em seguida abate 50 reais na conta 5000
      payOrder(self, 5000, self.orderID) # 110 - 50 = 60
      #adiciona mais um item com 4 unidades do produto
      createOrderItem(self, self.orderID, 4, self.currentRun) #60 + 40 = 100
      time.sleep(getRandomTime())
      #efetua o pagamento total do saldo devedor da ordem
      payOrder(self, 10000, self.orderID)
    else:
      print "no orderID for createOrder"

  @task(5)
  def getOrder(self):
    if (self.orderID != None):
      getOrder(self, self.orderID, self.currentRun)
    else:
      print "no orderID for getOrder"

class LocustInit(HttpLocust):
  task_set = RedWeddingBehaviour
  host = "http://52.87.172.243:80"
  #host = "http://sdlneuredelb-1416286239.us-east-1.elb.amazonaws.com:80"
