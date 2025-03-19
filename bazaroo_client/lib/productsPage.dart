import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'nav/customer_nav.dart';
import 'productDetails.dart';

class ProductsPage extends StatefulWidget {
  final String userId;
  ProductsPage({required this.userId});

  @override
  _ProductsPageState createState() => _ProductsPageState();
}

class _ProductsPageState extends State<ProductsPage> {
  List products = [];
  bool isLoading = true;
  bool hasError = false;

  @override
  void initState() {
    super.initState();
    fetchProducts();
  }

  Future<void> fetchProducts() async {
    final url = Uri.parse('http://localhost:3000/v1/products');
    try {
      final response = await http.get(url);
      if (response.statusCode == 200) {
        List jsonData = json.decode(response.body);
        setState(() {
          products = jsonData;
          isLoading = false;
        });
      } else {
        throw Exception('Failed to load products: ${response.statusCode}');
      }
    } catch (e) {
      setState(() {
        hasError = true;
        isLoading = false;
      });
      print('Error fetching products: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        toolbarHeight: 80,
        backgroundColor: Colors.red,
        title: const Text(
          'Products',
          style: TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold),
        ),
        automaticallyImplyLeading: false,
        actions: [
          Container(
            margin: const EdgeInsets.symmetric(horizontal: 25.0),
            child: IconButton(
              icon: const Icon(Icons.chevron_right, color: Colors.white, size: 35),
              onPressed: () {
                Navigator.pop(context);
              },
            ),
          ),
        ],
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(20.0),
          child: SizedBox(
              child: GridView.builder(
                shrinkWrap: true,
                physics: NeverScrollableScrollPhysics(),
                gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisSpacing: 16,
                  mainAxisSpacing: 16,
                  crossAxisCount: 2,
                  childAspectRatio: 0.75,
                ),
                itemCount: products.length,
                itemBuilder: (context, index) {
                  final product = products[index];
                  String imageUrl = "http://localhost:3000/${product['prod_image']['String']}";
                  return GestureDetector(
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => ProductDetailScreen(
                            userId: widget.userId,
                            prodId: product['prod_id'],
                          ),
                        ),
                      );
                    },
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        AspectRatio(
                          aspectRatio: 1,
                          child: Container(
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(15),
                              image: DecorationImage(
                                image: NetworkImage(imageUrl),
                                fit: BoxFit.cover,
                              ),
                            ),
                          ),
                        ),
                        SizedBox(height: 5),
                        Text(
                          product['prod_name'],
                          style: const TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                        Text(
                          "₱${product['buy_price']}",
                          style: const TextStyle(
                            fontSize: 16,
                            color: Colors.red,
                          ),
                        ),
                      ],
                    ),
                  );
                },
              ),
          ),
        ),
      ),
      bottomNavigationBar: BottomNavBar(userId: widget.userId),
    );
  }
}
