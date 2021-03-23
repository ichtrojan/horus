<template>
  <tr>
    <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
      <p class="text-gray-900 whitespace-no-wrap">{{ log.key }}</p>
    </td>
    <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
      <p class="text-gray-900 whitespace-no-wrap">/</p>
    </td>
    <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
        <span
            class="relative inline-block px-3 py-1 font-semibold text-green-900 leading-tight"
        >
          <span
              aria-hidden
              class="absolute inset-0 bg-green-200 opacity-50 rounded-full"
          ></span>
          <span class="relative">200</span>
        </span>
    </td>
    <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
      0.01ms
    </td>
    <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
      21st Jan, 2021
    </td>
    <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">View</td>
  </tr>
</template>

<script>
module.exports = {
  props: ['log'],
};
</script>
