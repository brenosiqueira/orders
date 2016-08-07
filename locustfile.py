from random import randint
from locust import HttpLocust, TaskSet, task, events
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
       "price": 200
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
       "unit_price": 10,
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

class RedWeddingBehaviour(TaskSet):

    @task
    def createOrder(self):
        currentRun = runCounter.next()
        #time.sleep(1)
        #Cria a Ordem
        startCreateOrder = getTimeInMillis()
        createOrderResponse = self.client.post("/orders", createOrderBody())
        endCreateOrder = getTimeInMillis() - startCreateOrder
        print "createOrder, responseCode[%d], time[%dms], run[%d]" % (createOrderResponse.status_code, endCreateOrder, currentRun)

        if (createOrderResponse.status_code == 200):
          #5 consultas simultaneas--
          orderID = createOrderResponse.json()['id']
          #Adiciona 5 itens de um produto a ordem, um OrderItem com o valor 5 na quantidade
          createOrderItem(self, orderID, 5, currentRun)
          #adiciona novamente um OrderItem com 5 unidades do produto
          createOrderItem(self, orderID, 5, currentRun)
          #adiciona um produto com somente uma unidade
          createOrderItem(self, orderID, 1, currentRun)
          #em seguida abate 50 reais na conta
          payOrder(self, 50, orderID)
          #adiciona mais um item com 4 unidades do produto
          createOrderItem(self, orderID, 4, currentRun)
          time.sleep(getRandomTime())
          #efetua o pagamento total do saldo devedor da ordem
          payOrder(self, 100, orderID)

class LocustInit(HttpLocust):
    task_set = RedWeddingBehaviour
    host = "http://localhost:80"
