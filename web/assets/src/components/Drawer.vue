<template>
    <mu-drawer :open="open" :docked="docked" @close="toggle()" zDepth="2">
      <mu-list @itemClick="docked ? '' : toggle()">
        <mu-list-item title="Status">
          <mu-icon slot="left" value="dashboard"/>
        </mu-list-item>
        <mu-list-item title="Counters">
          <mu-icon slot="left" value="plus_one"/>
        </mu-list-item>
        <mu-list-item v-if="docked" @click.native="open = false" title="Close"/>
      </mu-list>
    </mu-drawer> 
</template>

<script>
export default {
  data () {
    return {
      open: false,
      docked: true
    }
  },
  methods: {
    toggle (flag) {
      this.open = !this.open
      this.docked = !flag
    }
  }
}
</script>

<style>
.mu-drawer {
  padding-top: 48px;
  z-index: 65534 !important;
}
.mu-overlay {
  z-index: 65533 !important;
}
</style>
