<script setup>
import { ref, onMounted } from 'vue'
import { getRemain } from '@/apis/order';
import { getImageUrl } from '@/apis/image';


// 订单列表
const RemainList = ref([])
const total =ref(0)

const params = ref({
    page: 1,
    pageSize: 5
})

const getRemainList = async () => {
  const res = await getRemain(params.value)
  RemainList.value = res.result
  total.value = res.count
}

onMounted(() => getRemainList())

//页数切换
const pageChange = (page) => {
  params.value.page = page
  getRemainList()
}

</script>

<template>
  <div class="order-container">
      <div class="main-container">
        <div class="holder-container" v-if="RemainList.length === 0">
          <el-empty description="暂无订单数据" />
        </div>
        <div v-else>
          <!-- 订单列表 -->
          <div class="order-item" v-for="order in RemainList" :key="order.Id">
            <div class="head">
              <span>发布时间：{{ order.CreatTime }}</span>
            </div>
            <div class="body">
              <div class="column goods">
                <ul>
                  <li :key="order.Skus.Id">
                    <a class="image" href="javascript:;">
                      <img :src="getImageUrl(order.Skus.Image)" alt="" />
                    </a>
                    <div class="info">
                      <p class="name ellipsis-2">
                        {{ order.Skus.Name }}
                      </p>
                      <p class="attr ellipsis">
                        <span>{{ order.Skus.AttrsText }}</span>
                      </p>
                    </div>
                    <div class="price">¥{{ order.Skus.RealPay?.toFixed(2) }}</div>
                    <div class="count">x1</div>
                  </li>
                </ul>
              </div>
              <div class="column action">
                <RouterLink :to="`/detail/${order.Skus.Id}`">
                <p><a href="javascript:;">查看详情</a></p>
                </RouterLink>
              </div>
            </div>
          </div>
          <!-- 分页 -->
          <div class="pagination-container">
            <el-pagination :total="total" :page-size="params.pageSize" @current-change="pageChange" background layout="prev, pager, next" />
          </div>
        </div>
      </div>
  </div>

</template>

<style scoped lang="scss">
.order-container {
  padding: 10px 20px;

  .pagination-container {
    display: flex;
    justify-content: center;
  }

  .main-container {
    min-height: 500px;

    .holder-container {
      min-height: 500px;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }
}

.order-item {
  margin-bottom: 20px;
  border: 1px solid #f5f5f5;

  .head {
    height: 50px;
    line-height: 50px;
    background: #f5f5f5;
    padding: 0 20px;
    overflow: hidden;

    span {
      margin-right: 20px;

      &.down-time {
        margin-right: 0;
        float: right;

        i {
          vertical-align: middle;
          margin-right: 3px;
        }

        b {
          vertical-align: middle;
          font-weight: normal;
        }
      }
    }

    .del {
      margin-right: 0;
      float: right;
      color: #999;
    }
  }

  .body {
    display: flex;
    align-items: stretch;

    .column {
      border-left: 1px solid #f5f5f5;
      text-align: center;
      padding: 20px;

      >p {
        padding-top: 10px;
      }

      &:first-child {
        border-left: none;
      }

      &.goods {
        flex: 1;
        padding: 0;
        align-self: center;

        ul {
          li {
            border-bottom: 1px solid #f5f5f5;
            padding: 10px;
            display: flex;

            &:last-child {
              border-bottom: none;
            }

            .image {
              width: 70px;
              height: 70px;
              border: 1px solid #f5f5f5;
            }

            .info {
              width: 220px;
              text-align: left;
              padding: 0 10px;

              p {
                margin-bottom: 5px;

                &.name {
                  height: 38px;
                }

                &.attr {
                  color: #999;
                  font-size: 12px;

                  span {
                    margin-right: 5px;
                  }
                }
              }
            }

            .price {
              width: 100px;
            }

            .count {
              width: 80px;
            }
          }
        }
      }

      &.state {
        width: 120px;

        .green {
          color: $xtxColor;
        }
      }

      &.amount {
        width: 200px;

        .red {
          color: $priceColor;
        }
      }

      &.action {
        width: 140px;

        a {
          display: block;

          &:hover {
            color: $xtxColor;
          }
        }
      }
    }
  }
}
</style>