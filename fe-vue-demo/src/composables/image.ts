// 图片上传相关的公共逻辑：校验、压缩、体积格式化

/** 单张图片的文件大小上限（校验的是原文件，不是压缩后的数据） */
export const IMAGE_FILE_MAX_SIZE = 2 * 1024 * 1024

export const formatFileSize = (size: number) => {
    if (size < 1024) return `${size}B`
    if (size < 1024 * 1024) return `${Math.round(size / 1024)}KB`

    return `${(size / 1024 / 1024).toFixed(1)}MB`
}

export const isImageFile = (file: File) => file.type.startsWith('image/')

const isBlobUrl = (url?: string) => Boolean(url?.startsWith('blob:'))

/** 释放 el-upload 为预览生成的 blob 地址，避免内存泄漏 */
export const revokeBlobUrl = (url?: string) => {
    if (isBlobUrl(url) && url) URL.revokeObjectURL(url)
}

/**
 * 把图片压缩成 dataURL，避免直接存原图撑爆 localStorage
 * @param maxSize 压缩后长边的最长像素
 */
export const compressImage = (file: File, maxSize = 1024) => new Promise<string>((resolve, reject) => {
    const reader = new FileReader()

    reader.onload = () => {
        const image = new Image()

        image.onload = () => {
            try {
                const scale = Math.min(1, maxSize / Math.max(image.width, image.height))
                const canvas = document.createElement('canvas')
                canvas.width = Math.max(1, Math.round(image.width * scale))
                canvas.height = Math.max(1, Math.round(image.height * scale))

                const context = canvas.getContext('2d')
                if (!context) {
                    reject(new Error('无法创建画布'))
                    return
                }

                // 输出 jpeg 没有透明通道，先铺白底，避免透明 png 变成黑底
                context.fillStyle = '#ffffff'
                context.fillRect(0, 0, canvas.width, canvas.height)
                context.drawImage(image, 0, 0, canvas.width, canvas.height)

                resolve(canvas.toDataURL('image/jpeg', 0.8))
            } catch (error) {
                reject(error instanceof Error ? error : new Error('图片压缩失败'))
            }
        }

        image.onerror = () => reject(new Error('图片读取失败'))
        image.src = String(reader.result)
    }

    reader.onerror = () => reject(new Error('图片读取失败'))
    reader.readAsDataURL(file)
})
