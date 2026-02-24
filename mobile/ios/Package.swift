// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "OmnixiusiOS",
    platforms: [.iOS(.v16)],
    products: [
        .library(name: "OmnixiusiOS", targets: ["OmnixiusiOS"]),
    ],
    targets: [
        .target(name: "OmnixiusiOS", path: "Sources"),
    ]
)
