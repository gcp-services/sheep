<template>
  <div id="app" class="layout">
    <div class="header">
    <mu-appbar title="Sheep" zDepth="1">
      <mu-icon-button icon="menu" slot="left" @click="toggle(true)"/>
    </mu-appbar>
    <drawer ref="nav"/>
    </div>
    <div class="content">
      <router-view/>
    </div>
  </div>
</template>

<script>
import drawer from './components/Drawer'
export default {
  name: 'app',
  components: {
    'drawer': drawer
  },
  methods: {
    toggle (flag) {
      this.$refs.nav.toggle(flag)
    }
  }
}
</script>

<style>
.mu-appbar {
  position: fixed;
  height: 48px !important;
  z-index: 65535 !important;
}
.content {
  display: inline-block;
  padding-left: 2%;
  padding-right: 2%;
  margin-top: 48px;
}
</style>
