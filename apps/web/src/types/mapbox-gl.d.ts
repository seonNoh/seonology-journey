// Minimal ambient declaration so TypeScript knows 'mapbox-gl' is a valid module
// when it is not installed as a dependency. The real types would come from
// @types/mapbox-gl; we avoid pulling them in because mapbox-gl is loaded
// dynamically only when the Mapbox token is present.
declare module 'mapbox-gl' {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const mapboxgl: any
  export default mapboxgl
}
