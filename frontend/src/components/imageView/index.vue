<script setup>
import { getImageUrl } from '@/apis/image';
import { useMouseInElement } from '@vueuse/core';
import { watch, ref } from 'vue';
const props = defineProps({
    // 商品图片列表
    imageList: Array
})
// console.log(props.imageList)
// 图片列表
// const imageList = [
//   "https://yanxuan-item.nosdn.127.net/d917c92e663c5ed0bb577c7ded73e4ec.png",
//   "https://yanxuan-item.nosdn.127.net/e801b9572f0b0c02a52952b01adab967.jpg",
//   "https://yanxuan-item.nosdn.127.net/b52c447ad472d51adbdde1a83f550ac2.jpg",
//   "https://yanxuan-item.nosdn.127.net/f93243224dc37674dfca5874fe089c60.jpg",
//   "https://yanxuan-item.nosdn.127.net/f881cfe7de9a576aaeea6ee0d1d24823.jpg"
// ]
//小图切换大图显示 获取鼠标位置的小图
const activeIndex = ref(0)
const enterhandler = (i) =>{
    activeIndex.value = i
}
//放大镜
const layer = {
    width: 200,
    height: 200
}
const left = ref(0)
const top = ref(0)
const target = ref(null)

const positionX = ref(0)
const positionY = ref(0)
const {elementX, elementY, isOutside } = useMouseInElement(target)
watch([elementX, elementY, isOutside], () =>{

    if(isOutside.value)
      return
    if(elementX.value < layer.width*0.5){
        left.value = 0
    }else if(elementX.value > layer.width*1.5){
        left.value = layer.width
    }else{
        left.value = elementX.value - layer.width*0.5
    }

    if(elementY.value < layer.height*0.5){
        top.value = 0
    }else if(elementY.value > layer.height*1.5){
        top.value = layer.height
    }else{
        top.value = elementY.value - layer.height*0.5
    }

    positionX.value = -left.value*2
    positionY.value = -top.value*2
})
</script>


<template>
  <div class="goods-image" v-if="imageList">
    <!-- 左侧大图-->
    <div class="middle" ref="target">
      <img :src="getImageUrl(imageList[activeIndex])" alt="" style="object-fit: fill;"/>
      <!-- 蒙层小滑块 -->
      <div class="layer" v-show="!isOutside" :style="{ left: `${left}px`, top: `${top}px` }"></div>
    </div>
    <!-- 小图列表 -->
    <ul class="small" v-if="imageList.length > 0">
      <li v-for="(img, i) in imageList" :key="i" @mouseenter="enterhandler(i)" :class="{active:i === activeIndex}">
        <img :src="getImageUrl(img)" alt="" />
      </li>
    </ul>
    <!-- 放大镜大图 -->
    <div class="large" :style="[
      {
        backgroundImage: `url(${getImageUrl(imageList[activeIndex])})`,
        backgroundPositionX: `${positionX}px`,
        backgroundPositionY: `${positionY}px`,
      },
    ]" v-show="!isOutside"></div>
  </div>
</template>

<style scoped lang="scss">
.goods-image {
  width: 480px;
  height: 400px;
  position: relative;
  display: flex;

  .middle {
    width: 400px;
    height: 400px;
    background: #f5f5f5;
  }
img{
  width: 400px;
  height: 400px;
}
  .large {
    position: absolute;
    top: 0;
    left: 412px;
    width: 400px;
    height: 400px;
    z-index: 500;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    background-repeat: no-repeat;
    // 背景图:盒子的大小 = 2:1  将来控制背景图的移动来实现放大的效果查看 background-position
    background-size: 800px 800px;
    background-color: #f8f8f8;
  }

  .layer {
    width: 200px;
    height: 200px;
    background: rgba(0, 0, 0, 0.2);
    // 绝对定位 然后跟随咱们鼠标控制left和top属性就可以让滑块移动起来
    left: 0;
    top: 0;
    position: absolute;
  }

  .small {
    width: 80px;

    li {
      width: 68px;
      height: 68px;
      margin-left: 12px;
      margin-bottom: 15px;
      cursor: pointer;

      &:hover,
      &.active {
        border: 2px solid $xtxColor;
      }
    }
  }
}
</style>