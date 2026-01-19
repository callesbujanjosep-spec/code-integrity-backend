package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "version": "1.0.0"})
    })

    v1 := r.Group("/api/v1")
    {
        v1.POST("/register", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "treeId": "example-tree-id",
                "message": "Tree processing started",
            })
        })

        v1.GET("/tree/:id", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "treeId": c.Param("id"),
                "status": "completed",
                "rootHash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
                "files": []string{"README.md", "main.go"},
            })
        })

        v1.GET("/tree/:id/status", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "status": "completed",
                "progress": 1.0,
            })
        })

        v1.GET("/tree/:id/proof/:file", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "proof": []string{
                    "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
                    "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef123456789",
                },
                "leafHash": "0x7890abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456",
            })
        })

        v1.POST("/verify-offchain", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "verified": true,
                "steps": []string{
                    "Step 1: Calculate leaf hash = keccak256(filepath + content)",
                    "Step 2: Verify Merkle proof using provided root and proof array",
                },
            })
        })

        v1.GET("/stats", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "trees_generated": 42,
                "roots_registered": 42,
                "total_files_processed": 150,
            })
        })

        v1.GET("/contract-abi", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "abi": `[{"inputs":[{"internalType":"bytes32","name":"root","type":"bytes32"}],"name":"registerRoot","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"root","type":"bytes32"},{"internalType":"bytes32","name":"leaf","type":"bytes32"},{"internalType":"bytes32[]","name":"proof","type":"bytes32[]"}],"name":"verifyProof","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"pure","type":"function"}]`,
            })
        })
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}
