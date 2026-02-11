def calculate_area(radius):
    """
    Calculates the area of a circle given its radius. ✨
    """
    import math

    # Check if the radius is a positive number 🚀
    if radius < 0:
        # Raise an error if radius is negative
        raise ValueError("The radius cannot be negative. ✅")

    # Calculate the area using the formula: PI * r squared
    area = math.pi * (radius**2)

    # Return the calculated area to the caller
    return area


# Example usage of the calculate_area function
if __name__ == "__main__":
    print(f"Area: {calculate_area(5)}")  # Output the result
