<script setup>
import Readme from '../README.md'
import 'github-markdown-css'
import './resize.js'
</script>

<template>
  <div class="markdown-body container">
    <Readme />
  </div>
</template>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
}
</style>
