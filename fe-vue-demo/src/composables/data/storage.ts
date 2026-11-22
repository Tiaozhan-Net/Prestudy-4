const DATABASE_NAME = 'tenzors-blog'
const DATABASE_VERSION = 1
const STORE_NAME = 'app-data'

const toStorableValue = <T>(value: T): T => {
    return JSON.parse(JSON.stringify(value)) as T
}

const openDatabase = () => new Promise<IDBDatabase>((resolve, reject) => {
    const request = indexedDB.open(DATABASE_NAME, DATABASE_VERSION)

    request.onupgradeneeded = () => {
        request.result.createObjectStore(STORE_NAME)
    }
    request.onsuccess = () => resolve(request.result)
    request.onerror = () => reject(request.error)
})

export const getStoredValue = async <T>(key: string): Promise<T | null> => {
    const database = await openDatabase()

    return new Promise((resolve, reject) => {
        const request = database.transaction(STORE_NAME, 'readonly').objectStore(STORE_NAME).get(key)
        request.onsuccess = () => resolve((request.result as T | undefined) ?? null)
        request.onerror = () => reject(request.error)
    })
}

export const setStoredValue = async <T>(key: string, value: T): Promise<void> => {
    const database = await openDatabase()
    const storableValue = toStorableValue(value)

    return new Promise((resolve, reject) => {
        const request = database.transaction(STORE_NAME, 'readwrite').objectStore(STORE_NAME).put(storableValue, key)
        request.onsuccess = () => resolve()
        request.onerror = () => reject(request.error)
    })
}

export const removeStoredValue = async (key: string): Promise<void> => {
    const database = await openDatabase()

    return new Promise((resolve, reject) => {
        const request = database.transaction(STORE_NAME, 'readwrite').objectStore(STORE_NAME).delete(key)
        request.onsuccess = () => resolve()
        request.onerror = () => reject(request.error)
    })
}

export const migrateLocalStorage = async (keys: string[]) => {
    for (const key of keys) {
        const legacyValue = localStorage.getItem(key)
        if (legacyValue === null) continue

        try {
            const existingValue = await getStoredValue<string>(key)
            if (existingValue === null) {
                await setStoredValue(key, JSON.parse(legacyValue))
            }
        } catch {
            
        }
        localStorage.removeItem(key)
    }
}

export const initializeStorage = () => migrateLocalStorage(['users', 'current-user'])
